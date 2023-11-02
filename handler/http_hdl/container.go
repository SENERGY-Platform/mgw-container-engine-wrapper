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

package http_hdl

import (
	"fmt"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

const ctrIdParam = "c"

type containersQuery struct {
	Name  string   `form:"name"`
	State string   `form:"state"`
	Label []string `form:"label"`
}

type containerLogQuery struct {
	MaxLines int    `form:"max_lines"`
	Since    string `form:"since"`
	Until    string `form:"until"`
}

type deleteContainerQuery struct {
	Force bool `form:"force"`
}

func getContainersH(a lib.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := containersQuery{}
		if err := c.ShouldBindQuery(&query); err != nil {
			_ = c.Error(model.NewInvalidInputError(err))
			return
		}
		filter := model.ContainerFilter{Name: query.Name}
		if query.State != "" {
			_, ok := model.ContainerStateMap[query.State]
			if !ok {
				_ = c.Error(model.NewInvalidInputError(fmt.Errorf("unknown container state '%s'", query.State)))
				return
			}
			filter.State = query.State
		}
		filter.Labels = GenLabels(query.Label)
		containers, err := a.GetContainers(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, containers)
	}
}

func postContainerH(a lib.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		container := model.Container{}
		if err := c.ShouldBindJSON(&container); err != nil {
			_ = c.Error(model.NewInvalidInputError(err))
			return
		}
		id, err := a.CreateContainer(c.Request.Context(), container)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.String(http.StatusOK, id)
	}
}

func deleteContainerH(a lib.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := deleteContainerQuery{}
		if err := c.ShouldBindQuery(&query); err != nil {
			_ = c.Error(model.NewInvalidInputError(err))
			return
		}
		if err := a.RemoveContainer(c.Request.Context(), c.Param(ctrIdParam), query.Force); err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func getContainerH(a lib.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		container, err := a.GetContainer(c.Request.Context(), c.Param(ctrIdParam))
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, container)
	}
}

func getContainerLogH(a lib.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		rc, err := a.GetContainerLog(c.Request.Context(), c.Param(ctrIdParam), logOptions)
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

func patchContainerStartH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		err := a.StartContainer(gc.Request.Context(), gc.Param(ctrIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func patchContainerStopH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		jID, err := a.StopContainer(gc.Request.Context(), gc.Param(ctrIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, jID)
	}
}

func patchContainerRestartH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		jID, err := a.RestartContainer(gc.Request.Context(), gc.Param(ctrIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, jID)
	}
}
