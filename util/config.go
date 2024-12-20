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

package util

import (
	"github.com/SENERGY-Platform/go-service-base/config-hdl"
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	envldr "github.com/y-du/go-env-loader"
	"github.com/y-du/go-log-level/level"
	"io/fs"
	"os"
	"reflect"
)

type JobsConfig struct {
	BufferSize  int   `json:"buffer_size" env_var:"JOBS_BUFFER_SIZE"`
	MaxNumber   int   `json:"max_number" env_var:"JOBS_MAX_NUMBER"`
	CCHInterval int   `json:"cch_interval" env_var:"JOBS_CCH_INTERVAL"`
	JHInterval  int   `json:"jh_interval" env_var:"JOBS_JH_INTERVAL"`
	PJHInterval int64 `json:"pjh_interval" env_var:"JOBS_PJH_INTERVAL"`
	MaxAge      int64 `json:"max_age" env_var:"JOBS_MAX_AGE"`
}

type SocketConfig struct {
	Path     string      `json:"path" env_var:"SOCKET_PATH"`
	GroupID  int         `json:"group_id" env_var:"SOCKET_GROUP_ID"`
	FileMode fs.FileMode `json:"file_mode" env_var:"SOCKET_FILE_MODE"`
}

type LoggerConfig struct {
	Level        level.Level `json:"level" env_var:"LOGGER_LEVEL"`
	Utc          bool        `json:"utc" env_var:"LOGGER_UTC"`
	Path         string      `json:"path" env_var:"LOGGER_PATH"`
	FileName     string      `json:"file_name" env_var:"LOGGER_FILE_NAME"`
	Terminal     bool        `json:"terminal" env_var:"LOGGER_TERMINAL"`
	Microseconds bool        `json:"microseconds" env_var:"LOGGER_MICROSECONDS"`
	Prefix       string      `json:"prefix" env_var:"LOGGER_PREFIX"`
}

type DockerConfig struct {
	Host          string `json:"host" env_var:"DOCKER_HOST"`
	CtrLogDriver  string `json:"ctr_log_driver" env_var:"DOCKER_CTR_LOG_DRIVER"`
	CtrLogMaxSize string `json:"ctr_log_max_size" env_var:"DOCKER_CTR_LOG_MAX_SIZE"`
	CtrLogMaxFile int    `json:"ctr_log_max_file" env_var:"DOCKER_CTR_LOG_MAX_FILE"`
}

type Config struct {
	Logger LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
	Socket SocketConfig `json:"socket" env_var:"SOCKET_CONFIG"`
	Jobs   JobsConfig   `json:"jobs" env_var:"JOBS_CONFIG"`
	Docker DockerConfig `json:"docker" env_var:"DOCKER_CONFIG"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		Logger: LoggerConfig{
			Level:        level.Warning,
			Utc:          true,
			Path:         "./",
			FileName:     "mgw_ce_wrapper",
			Microseconds: true,
		},
		Socket: SocketConfig{
			Path:     "./ce_wrapper.sock",
			GroupID:  os.Getgid(),
			FileMode: 0660,
		},
		Jobs: JobsConfig{
			BufferSize:  200,
			MaxNumber:   20,
			CCHInterval: 500000,
			JHInterval:  500000,
			PJHInterval: 300000000000,
			MaxAge:      172800000000000,
		},
		Docker: DockerConfig{
			Host: "unix:///var/run/docker.sock",
		},
	}
	err := config_hdl.Load(&cfg, nil, map[reflect.Type]envldr.Parser{reflect.TypeOf(level.Off): sb_logger.LevelParser}, nil, path)
	return &cfg, err
}
