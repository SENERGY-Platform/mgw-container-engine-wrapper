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

package engine

import (
	"deployment-manager/util/logger"
	"github.com/gin-gonic/gin"
)

type Config struct {
	StaticOrigins []string `json:"static_origins" env_var:"API_ENGINE_STATIC_ORIGINS"`
}

func New(config Config) *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.ForwardedByClientIP = false
	e.Use(ginLogger(gin.LoggerConfig{}, logger.GetLevel()), gin.Recovery(), checkStaticOrigin(config.StaticOrigins))
	return e
}
