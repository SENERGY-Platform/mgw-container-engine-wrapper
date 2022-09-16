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
	"deployment-manager/manager/handler/gin-web"
	"encoding/json"
	envldr "github.com/y-du/go-env-loader"
	"github.com/y-du/go-log-level/level"
	"os"
	"reflect"
)

type Config struct {
	SocketPath string         `json:"socket_path" env_var:"SOCKET_PATH"`
	Logger     LoggerConfig   `json:"logger" env_var:"LOG_CONFIG"`
	ApiEngine  gin_web.Config `json:"api_engine" env_var:"API_ENGINE_CONFIG"`
}

func NewConfig(path *string) (cfg *Config, err error) {
	cfg = &Config{
		SocketPath: "/opt/deployment-manager/manager.sock",
		Logger: LoggerConfig{
			Level:    level.Warning,
			Utc:      true,
			Prefix:   "[DM] ",
			Path:     "/var/log/",
			FileName: "mgw-deployment-manager",
		},
	}
	if path != nil {
		var file *os.File
		if file, err = os.Open(*path); err != nil {
			return
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(cfg); err != nil {
			return
		}
	}
	err = envldr.LoadEnvUserParser(cfg, nil, typeParsers, nil)
	return cfg, err
}

var typeParsers = map[reflect.Type]envldr.Parser{
	reflect.TypeOf(level.Off): LogLevelParser,
}

var LogLevelParser envldr.Parser = func(t reflect.Type, val string, params []string, kwParams map[string]string) (interface{}, error) {
	return level.Parse(val)
}
