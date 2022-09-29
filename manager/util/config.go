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

package util

import (
	"github.com/SENERGY-Platform/go-service-base"
	"github.com/y-du/go-log-level/level"
)

type Config struct {
	SocketPath string                `json:"socket_path" env_var:"SOCKET_PATH"`
	Logger     srv_base.LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
}

func NewConfig(path *string) (*Config, error) {
	cfg := Config{
		SocketPath: "/opt/deployment-manager/manager.sock",
		Logger: srv_base.LoggerConfig{
			Level:        level.Warning,
			Utc:          true,
			Path:         "/var/log/",
			FileName:     "mgw-deployment-manager",
			Microseconds: true,
		},
	}
	err := srv_base.LoadConfig(path, &cfg)
	return &cfg, err
}
