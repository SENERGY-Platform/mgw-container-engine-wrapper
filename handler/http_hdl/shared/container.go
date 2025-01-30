/*
 * Copyright 2025 InfAI (CC SES)
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

package shared

import (
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path"
	"time"
)

type containerLogQuery struct {
	MaxLines int    `form:"max_lines"`
	Since    string `form:"since"`
	Until    string `form:"until"`
}

// getContainerLogH godoc
// @Summary Get container log
// @Description Get a container's log.
// @Tags Containers
// @Produce	plain
// @Param id path string true "container ID"
// @Param max_lines query integer false "max num of lines"
// @Param since query string false "RFC3339Nano timestamp"
// @Param until query string false "RFC3339Nano timestamp"
// @Success	200 {string} string "log"
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /logs/{id} [get]
func getContainerLogH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(model.ContainerLogsPath, ":id"), func(c *gin.Context) {
		query := containerLogQuery{}
		if err := c.ShouldBindQuery(&query); err != nil {
			_ = c.Error(model.NewInvalidInputError(err))
			return
		}
		logOptions := model.LogFilter{MaxLines: query.MaxLines}
		if query.Since != "" {
			t, err := time.Parse(time.RFC3339Nano, query.Since)
			if err != nil {
				_ = c.Error(model.NewInvalidInputError(err))
				return
			}
			logOptions.Since = t
		}
		if query.Until != "" {
			t, err := time.Parse(time.RFC3339Nano, query.Until)
			if err != nil {
				_ = c.Error(model.NewInvalidInputError(err))
				return
			}
			logOptions.Until = t
		}
		rc, err := a.GetContainerLog(c.Request.Context(), c.Param("id"), logOptions)
		if err != nil {
			_ = c.Error(err)
			return
		}
		defer rc.Close()
		c.Status(http.StatusOK)
		c.Header("Transfer-Encoding", "chunked")
		c.Header("Content-Type", gin.MIMEPlain)
		for {
			var b = make([]byte, 204800)
			n, rErr := rc.Read(b)
			if rErr != nil {
				if rErr == io.EOF {
					if n > 0 {
						_, wErr := c.Writer.Write(b[:n])
						if wErr != nil {
							_ = c.Error(model.NewInternalError(wErr))
							return
						}
						c.Writer.Flush()
					}
					break
				}
				_ = c.Error(model.NewInternalError(rErr))
				return
			}
			_, wErr := c.Writer.Write(b)
			if wErr != nil {
				_ = c.Error(model.NewInternalError(wErr))
				return
			}
			c.Writer.Flush()
		}
	}
}
