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
	"deployment-manager/manager/handler"
	"deployment-manager/manager/handler/api-engine"
	"deployment-manager/util"
	"deployment-manager/util/configuration"
	"deployment-manager/util/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flags := util.NewFlags()
	config, err := configuration.NewConfig(flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger.InitLogger(config.Logger.Level, "[DM] ", false, config.Logger.Utc)
	engine := api_engine.NewEngine(config.StaticOrigins, config.Logger.Level)
	api.SetRoutes(engine)
	server := http.Server{
		Handler: engine,
	}
	listener, err := handler.NewUnixListener(config.SocketPath)
	if err != nil {
		logger.Fatal(err)
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-shutdown
		logger.WarningF("received signal '%s'", sig)
		logger.Info("initiating shutdown ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server forced to shutdown: ", err)
		}
	}()
	logger.Info("starting server ...")
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		logger.Fatal("starting server failed: ", err)
	} else {
		logger.Info("shutdown complete")
	}
}
