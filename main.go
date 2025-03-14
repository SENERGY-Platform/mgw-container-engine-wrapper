/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/go-cc-job-handler/ccjh"
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker_hdl"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/http_hdl"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/wrapper"
	"github.com/SENERGY-Platform/mgw-go-service-base/job-hdl"
	"github.com/SENERGY-Platform/mgw-go-service-base/srv-info-hdl"
	sb_util "github.com/SENERGY-Platform/mgw-go-service-base/util"
	"github.com/SENERGY-Platform/mgw-go-service-base/watchdog"
	"github.com/docker/docker/client"
	"net/http"
	"os"
	"syscall"
	"time"
)

var version string

func main() {
	srvInfoHdl := srv_info_hdl.New("ce-wrapper", version)

	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	util.ParseFlags()

	config, err := util.NewConfig(util.Flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ec = 1
		return
	}

	logFile, err := util.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *sb_logger.LogFileError
		if errors.As(err, &logFileError) {
			ec = 1
			return
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	util.Logger.Printf("%s %s", srvInfoHdl.GetName(), srvInfoHdl.GetVersion())

	util.Logger.Debugf("config: %s", sb_util.ToJsonStr(config))

	watchdog.Logger = util.Logger
	wtchdg := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	dockerClient, err := client.NewClientWithOpts(client.WithTLSClientConfigFromEnv(), client.WithHost(config.Docker.Host), client.WithVersionFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	defer dockerClient.Close()

	dockerHandler, err := docker_hdl.New(dockerClient, docker_hdl.ContainerLogConf{
		Driver:  config.Docker.CtrLogDriver,
		MaxSize: config.Docker.CtrLogMaxSize,
		MaxFile: config.Docker.CtrLogMaxFile,
	})
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	ccHandler := ccjh.New(config.Jobs.BufferSize)

	job_hdl.Logger = util.Logger
	job_hdl.ErrCodeMapper = util.GetErrCode
	job_hdl.NewNotFoundErr = model.NewNotFoundError
	job_hdl.NewInvalidInputError = model.NewInvalidInputError
	job_hdl.NewInternalErr = model.NewInternalError
	jobCtx, jobCF := context.WithCancel(context.Background())
	jobHandler := job_hdl.New(jobCtx, ccHandler)
	purgeJobsHdl := job_hdl.NewPurgeJobsHandler(jobHandler, time.Duration(config.Jobs.PJHInterval), time.Duration(config.Jobs.MaxAge))

	wtchdg.RegisterStopFunc(func() error {
		ccHandler.Stop()
		jobCF()
		if ccHandler.Active() > 0 {
			util.Logger.Info("waiting for active jobs to cancel ...")
			ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
			defer cf()
			for ccHandler.Active() != 0 {
				select {
				case <-ctx.Done():
					return fmt.Errorf("canceling jobs took too long")
				default:
					time.Sleep(50 * time.Millisecond)
				}
			}
			util.Logger.Info("jobs canceled")
		}
		return nil
	})

	cew := wrapper.New(dockerHandler, jobHandler, srvInfoHdl)

	httpHandler, err := http_hdl.New(cew, map[string]string{
		model.HeaderApiVer:  srvInfoHdl.GetVersion(),
		model.HeaderSrvName: srvInfoHdl.GetName(),
	})
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	listener, err := sb_util.NewUnixListener(config.Socket.Path, os.Getuid(), config.Socket.GroupID, config.Socket.FileMode)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	server := &http.Server{Handler: httpHandler}
	srvCtx, srvCF := context.WithCancel(context.Background())
	wtchdg.RegisterStopFunc(func() error {
		if srvCtx.Err() == nil {
			ctxWt, cf := context.WithTimeout(context.Background(), time.Second*5)
			defer cf()
			if err := server.Shutdown(ctxWt); err != nil {
				return err
			}
			util.Logger.Info("http server shutdown complete")
		}
		return nil
	})
	wtchdg.RegisterHealthFunc(func() bool {
		if srvCtx.Err() == nil {
			return true
		}
		util.Logger.Error("http server closed unexpectedly")
		return false
	})

	wtchdg.Start()

	err = ccHandler.RunAsync(config.Jobs.MaxNumber, time.Duration(config.Jobs.JHInterval*1000))
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	purgeJobsHdl.Start(jobCtx)

	dCtx, dCF := context.WithCancel(context.Background())
	wtchdg.RegisterStopFunc(func() error {
		dCF()
		return nil
	})

	go func() {
		defer dCF()
		dockerInfo, err := dockerHandler.ServerInfo(dCtx, time.Millisecond*100)
		if err != nil {
			util.Logger.Error(err)
			ec = 1
			wtchdg.Trigger()
			return
		}
		util.Logger.Debugf("docker: %s", sb_util.ToJsonStr(dockerInfo))
	}()

	go func() {
		defer srvCF()
		util.Logger.Info("starting http server ...")
		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			util.Logger.Error(err)
			ec = 1
			return
		}
	}()

	ec = wtchdg.Join()
}
