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

package gin_web

import (
	"github.com/gin-gonic/gin"
	"github.com/y-du/go-log-level/level"
	"os"
)

func New(logLevel level.Level, logFile *os.File) *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.ForwardedByClientIP = false
	lc := gin.LoggerConfig{}
	if logFile != nil {
		lc.Output = logFile
	}
	e.Use(ginLogger(lc, logLevel), gin.Recovery())
	return e
}
