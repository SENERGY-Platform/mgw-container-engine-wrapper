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
	"deployment-manager/manager/api"
	"deployment-manager/manager/handler"
	"deployment-manager/manager/handler/gin-engine"
	"deployment-manager/util"
	"deployment-manager/util/configuration"
	"deployment-manager/util/logger"
	"fmt"
	"os"
)

func main() {
	flags := util.NewFlags()
	config, err := configuration.NewConfig(flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger.InitLogger(config.Logger.Level, "[DM] ", false, config.Logger.Utc)
	listener, err := handler.NewUnixListener(config.SocketPath)
	if err != nil {
		logger.Fatal(err)
	}
	engine := gin_engine.NewEngine(config.StaticOrigins, config.Logger.Level)
	api.SetRoutes(engine)

}
