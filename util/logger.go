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
	envldr "github.com/y-du/go-env-loader"
	log_level "github.com/y-du/go-log-level"
	"github.com/y-du/go-log-level/level"
	"log"
	"os"
	"reflect"
)

type LoggerConfig struct {
	Level  level.Level `json:"level" env_var:"LOGGER_LEVEL"`
	Utc    bool        `json:"utc" env_var:"LOGGER_UTC"`
	Prefix string      `json:"prefix" env_var:"LOGGER_PREFIX"`
}

var Logger *log_level.Logger

func InitLogger(config LoggerConfig) (err error) {
	flags := log.Ldate | log.Ltime | log.Lmsgprefix
	if config.Utc {
		flags = flags | log.LUTC
	}
	Logger, err = log_level.New(log.New(os.Stderr, config.Prefix, flags), config.Level)
	return
}

var LogLevelParser envldr.Parser = func(t reflect.Type, val string, params []string, kwParams map[string]string) (interface{}, error) {
	return level.Parse(val)
}
