/*
 * Copyright 2022 InfAI (CC SES)
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
	"deployment-manager/manager/api"
	"deployment-manager/manager/handler/docker"
	"deployment-manager/manager/util"
	"deployment-manager/manager/util/gin-mw"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/go-service-base"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version string

func main() {
	srv_base.PrintInfo("mgw-deployment-manager", version)

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
	dockerHandler := docker.New(dockerClient)
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}
	defer dockerHandler.Close()
	dockerInfo, err := dockerHandler.ServerInfo(context.Background())
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}
	srv_base.Logger.Debugf("docker: %s", srv_base.ToJsonStr(dockerInfo))

	gin.SetMode(gin.ReleaseMode)
	apiEngine := gin.New()
	apiEngine.Use(gin_mw.LoggerHandler(srv_base.Logger), gin_mw.ErrorHandler, gin.Recovery())
	dmApi := api.New(dockerHandler)
	dmApi.SetRoutes(apiEngine)

	listener, err := util.NewUnixListener(config.SocketPath)
	if err != nil {
		srv_base.Logger.Error(err)
		return
	}
	server := http.Server{
		Handler: apiEngine,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-shutdown
		srv_base.Logger.Warningf("received signal '%s'", sig)
		srv_base.Logger.Info("initiating shutdown ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			srv_base.Logger.Error("server forced to shutdown: ", err)
		}
	}()

	srv_base.Logger.Info("starting server ...")
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		srv_base.Logger.Error("starting server failed: ", err)
		return
	} else {
		srv_base.Logger.Info("shutdown complete")
	}
}
