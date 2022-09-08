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

package configuration

import (
	"deployment-manager/util/logger"
)

type LoggerConfig struct {
	Level logger.Level `json:"level" env_var:"LOG_LEVEL"`
	Utc   bool         `json:"utc" env_var:"LOG_UTC"`
}

type Config struct {
	SocketPath    string       `json:"socket_path" env_var:"SOCKET_PATH"`
	StaticOrigins []string     `json:"static_origins" env_var:"STATIC_ORIGINS"`
	Logger        LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
}
