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
	"fmt"
	log_level "github.com/y-du/go-log-level"
	"github.com/y-du/go-log-level/level"
	"log"
	"os"
	"path"
)

var Logger *log_level.Logger

type LoggerConfig struct {
	Level        level.Level `json:"level" env_var:"LOGGER_LEVEL"`
	Utc          bool        `json:"utc" env_var:"LOGGER_UTC"`
	Path         string      `json:"path" env_var:"LOGGER_PATH"`
	FileName     string      `json:"file_name" env_var:"LOGGER_FILE_NAME"`
	Terminal     bool        `json:"terminal" env_var:"LOGGER_TERMINAL"`
	Microseconds bool        `json:"microseconds" env_var:"LOGGER_MICROSECONDS"`
	Prefix       string      `json:"prefix" env_var:"LOGGER_PREFIX"`
}

type LogFileError struct {
	msg string
}

func (e LogFileError) Error() string {
	return e.msg
}

func InitLogger(config LoggerConfig) (out *os.File, err error) {
	flags := log.Ldate | log.Ltime | log.Lmsgprefix
	if config.Utc {
		flags = flags | log.LUTC
	}
	if config.Microseconds {
		flags = flags | log.Lmicroseconds
	}
	if config.Terminal {
		out = os.Stderr
	} else {
		out, err = os.OpenFile(path.Join(config.Path, fmt.Sprintf("%s.log", config.FileName)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			err = &LogFileError{
				msg: err.Error(),
			}
			return
		}
	}
	Logger, err = log_level.New(log.New(out, config.Prefix, flags), config.Level)
	return
}
