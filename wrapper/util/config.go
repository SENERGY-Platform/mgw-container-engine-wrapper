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
	"github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/y-du/go-log-level/level"
	"io/fs"
	"os"
)

type HeaderConfig struct {
	RequestUri string `json:"request_uri" env_var:"HEADER_REQUEST_URI"`
	Uri        string `json:"uri" env_var:"HEADER_URI"`
}

type JobsConfig struct {
	BufferSize int `json:"buffer_size" env_var:"JOBS_BUFFER_SIZE"`
	MaxNumber  int `json:"max_number" env_var:"JOBS_MAX_NUMBER"`
}

type SocketConfig struct {
	Path     string      `json:"path" env_var:"SOCKET_PATH"`
	GroupID  int         `json:"group_id" env_var:"SOCKET_GROUP_ID"`
	FileMode fs.FileMode `json:"file_mode" env_var:"SOCKET_FILE_MODE"`
}

type Config struct {
	Logger srv_base.LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
	Socket SocketConfig          `json:"socket" env_var:"SOCKET_CONFIG"`
	Jobs   JobsConfig            `json:"jobs" env_var:"JOBS_CONFIG"`
	Header HeaderConfig          `json:"header" env_var:"HEADER_CONFIG"`
}

func NewConfig(path *string) (*Config, error) {
	cfg := Config{
		Logger: srv_base.LoggerConfig{
			Level:        level.Warning,
			Utc:          true,
			Path:         "/var/log/",
			FileName:     "mgw-container-engine-wrapper",
			Microseconds: true,
		},
		Socket: SocketConfig{
			Path:     "/opt/mgw/sockets/cew.sock",
			GroupID:  os.Getgid(),
			FileMode: 0660,
		},
		Jobs: JobsConfig{
			BufferSize: 50,
			MaxNumber:  10,
		},
		Header: HeaderConfig{
			RequestUri: "X-Request-Uri",
			Uri:        "X-Uri",
		},
	}
	err := srv_base.LoadConfig(path, &cfg, nil, nil, nil)
	return &cfg, err
}
