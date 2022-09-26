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

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/y-du/go-log-level/level"
	"os"
	"time"
)

func logFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("%v %3d | %v | %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		param.StatusCode,
		param.Latency,
		//param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

func Logger(conf gin.LoggerConfig, lvl level.Level) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = logFormatter
	}
	out := conf.Output
	if out == nil {
		out = os.Stderr
	}
	notlogged := conf.SkipPaths
	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}
			param.TimeStamp = time.Now().UTC()
			param.Latency = param.TimeStamp.Sub(start)
			//param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			if lvl > level.Info {
				fmt.Fprint(out, formatter(param))
			} else if lvl <= level.Info && len(c.Errors) > 0 {
				fmt.Fprint(out, formatter(param))
			}
		}
	}
}
