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
	"deployment-manager/manager/handler/gin-web"
	"deployment-manager/util"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version string

func main() {
	util.PrintInfo("mgw-deployment-manager", version)

	flags := util.NewFlags()
	config, err := util.NewConfig(flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logFile, err := util.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *util.LogFileError
		if errors.As(err, &logFileError) {
			os.Exit(1)
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	util.Logger.Debugf("config: %s", util.ToJsonStr(config))

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		util.Logger.Error(err)
		return
	}
	dockerHandler := docker.New(dockerClient)
	if err != nil {
		util.Logger.Error(err)
		return
	}
	defer dockerHandler.Close()
	dockerInfo, err := dockerHandler.ServerInfo(context.Background())
	if err != nil {
		util.Logger.Error(err)
		return
	}
	util.Logger.Debugf("docker: %s", util.ToJsonStr(dockerInfo))

	dmApi := api.New(dockerHandler)
	apiEngine := gin_web.New(config.ApiEngine, config.Logger.Level, logFile)
	api.SetRoutes(apiEngine, dmApi)

	listener, err := util.NewUnixListener(config.SocketPath)
	if err != nil {
		util.Logger.Error(err)
		return
	}
	server := http.Server{
		Handler: apiEngine,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-shutdown
		util.Logger.Warningf("received signal '%s'", sig)
		util.Logger.Info("initiating shutdown ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			util.Logger.Error("server forced to shutdown: ", err)
		}
	}()

	util.Logger.Info("starting server ...")
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		util.Logger.Error("starting server failed: ", err)
		return
	} else {
		util.Logger.Info("shutdown complete")
	}
}
