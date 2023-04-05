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
	"github.com/SENERGY-Platform/go-service-base/srv-base/types"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/api"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/api/http_engine"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/api/http_engine/middleware"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/job"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var version string

func main() {
	srv_base.PrintInfo("mgw-container-engine-wrapper", version)

	flags := util.NewFlags()

	config, err := util.NewConfig(flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logFile, err := srv_base.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *srv_base.LogFileError
		if errors.As(err, &logFileError) {
			os.Exit(1)
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	srv_base.Logger.Debugf("config: %s", srv_base.ToJsonStr(config))

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}
	defer dockerClient.Close()

	dockerHandler := docker.New(dockerClient)

	dockerInfo, err := dockerHandler.ServerInfo(context.Background())
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}
	srv_base.Logger.Debugf("docker: %s", srv_base.ToJsonStr(dockerInfo))

	ccHandler := ccjh.New(config.Jobs.BufferSize)

	jobCtx, cFunc := context.WithCancel(context.Background())
	jobHandler := job.New(jobCtx, ccHandler)

	defer func() {
		ccHandler.Stop()
		cFunc()
		if ccHandler.Active() > 0 {
			srv_base.Logger.Info("waiting for active jobs to cancel ...")
			ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
			defer cf()
			for ccHandler.Active() != 0 {
				select {
				case <-ctx.Done():
					srv_base.Logger.Error("canceling jobs took too long")
					return
				default:
					time.Sleep(50 * time.Millisecond)
				}
			}
			srv_base.Logger.Info("jobs canceled")
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	apiEngine := gin.New()
	apiEngine.Use(gin_mw.LoggerHandler(srv_base.Logger), middleware.ErrorHandler, gin.Recovery())
	apiEngine.UseRawPath = true
	cewApi := api.New(dockerHandler, jobHandler)

	http_engine.SetRoutes(apiEngine, cewApi)

	listener, err := srv_base.NewUnixListener(config.Socket.Path, os.Getuid(), config.Socket.GroupID, config.Socket.FileMode)
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}

	err = ccHandler.RunAsync(config.Jobs.MaxNumber, time.Duration(config.Jobs.JHInterval*1000))
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}

	srv_base.StartServer(&http.Server{Handler: apiEngine}, listener, srv_base_types.DefaultShutdownSignals)
}
