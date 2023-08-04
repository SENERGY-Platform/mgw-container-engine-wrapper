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
	"github.com/SENERGY-Platform/gin-middleware"
	"github.com/SENERGY-Platform/go-cc-job-handler/ccjh"
	"github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/api"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker_hdl"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/http_hdl"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/job_hdl"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/docker/docker/client"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"syscall"
	"time"
)

var version string

func main() {
	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	srv_base.PrintInfo(model.ServiceName, version)

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
		var logFileError *srv_base.LogFileError
		if errors.As(err, &logFileError) {
			ec = 1
			return
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	util.Logger.Debugf("config: %s", srv_base.ToJsonStr(config))

	watchdog := srv_base.NewWatchdog(util.Logger, syscall.SIGINT, syscall.SIGTERM)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	defer dockerClient.Close()

	dockerHandler := docker_hdl.New(dockerClient)

	ccHandler := ccjh.New(config.Jobs.BufferSize)

	jobCtx, jobCF := context.WithCancel(context.Background())
	jobHandler := job_hdl.New(jobCtx, ccHandler)

	watchdog.RegisterStopFunc(func() error {
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

	gin.SetMode(gin.ReleaseMode)
	httpHandler := gin.New()
	staticHeader := map[string]string{
		model.HeaderApiVer:  version,
		model.HeaderSrvName: model.ServiceName,
	}
	httpHandler.Use(gin_mw.StaticHeaderHandler(staticHeader), requestid.New(requestid.WithCustomHeaderStrKey(model.HeaderRequestID)), gin_mw.LoggerHandler(util.Logger, nil, func(gc *gin.Context) string {
		return requestid.Get(gc)
	}), gin_mw.ErrorHandler(http_hdl.GetStatusCode, ", "), gin.Recovery())
	httpHandler.UseRawPath = true
	cewApi := api.New(dockerHandler, jobHandler)

	http_hdl.SetRoutes(httpHandler, cewApi)
	util.Logger.Debugf("routes: %s", srv_base.ToJsonStr(http_hdl.GetRoutes(httpHandler)))

	listener, err := srv_base.NewUnixListener(config.Socket.Path, os.Getuid(), config.Socket.GroupID, config.Socket.FileMode)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	server := &http.Server{Handler: httpHandler}
	srvCtx, srvCF := context.WithCancel(context.Background())
	watchdog.RegisterStopFunc(func() error {
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
	watchdog.RegisterHealthFunc(func() bool {
		if srvCtx.Err() == nil {
			return true
		}
		util.Logger.Error("http server closed unexpectedly")
		return false
	})

	watchdog.Start()

	err = ccHandler.RunAsync(config.Jobs.MaxNumber, time.Duration(config.Jobs.JHInterval*1000))
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	dCtx, dCF := context.WithCancel(context.Background())
	watchdog.RegisterStopFunc(func() error {
		dCF()
		return nil
	})

	go func() {
		defer dCF()
		dockerInfo, err := dockerHandler.ServerInfo(dCtx, time.Millisecond*100)
		if err != nil {
			util.Logger.Error(err)
			ec = 1
			watchdog.Trigger()
			return
		}
		util.Logger.Debugf("docker: %s", srv_base.ToJsonStr(dockerInfo))
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

	ec = watchdog.Join()
}
