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

package http_engine

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
	MaxLines int   `form:"max_lines"`
	Since    int64 `form:"since"`
	Until    int64 `form:"until"`
}

func getContainersH(a lib.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := containersQuery{}
		if err := c.ShouldBindQuery(&query); err != nil {
			c.Status(http.StatusBadRequest)
			_ = c.Error(err)
			return
		}
		filter := model.ContainerFilter{Name: query.Name}
		if query.State != "" {
			_, ok := model.ContainerStateMap[query.State]
			if !ok {
				c.Status(http.StatusBadRequest)
				_ = c.Error(fmt.Errorf("unknown container state '%s'", query.State))
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
			c.Status(http.StatusBadRequest)
			_ = c.Error(err)
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
		if err := a.RemoveContainer(c.Request.Context(), c.Param(ctrIdParam)); err != nil {
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
			c.Status(http.StatusBadRequest)
			_ = c.Error(err)
			return
		}
		logOptions := model.LogFilter{MaxLines: query.MaxLines}
		if query.Since > 0 {
			logOptions.Since = time.UnixMicro(query.Since)
		}
		if query.Until > 0 {
			logOptions.Until = time.UnixMicro(query.Until)
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
							c.Status(http.StatusInternalServerError)
							_ = c.Error(wErr)
							return
						}
						c.Writer.Flush()
					}
					break
				}
				c.Status(http.StatusInternalServerError)
				_ = c.Error(rErr)
				return
			}
			_, wErr := c.Writer.Write(b)
			if wErr != nil {
				c.Status(http.StatusInternalServerError)
				_ = c.Error(wErr)
				return
			}
			c.Writer.Flush()
		}
	}
}

func postContainerCtrlH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var ctrlReq model.ContainerCtrlRequest
		if err := gc.ShouldBindJSON(&ctrlReq); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		switch ctrlReq.State {
		case model.RunningState:
			err := a.StartContainer(gc.Request.Context(), gc.Param(ctrIdParam))
			if err != nil {
				_ = gc.Error(err)
				return
			}
			gc.Status(http.StatusOK)
		case model.StoppedState:
			id, err := a.StopContainer(gc.Request.Context(), gc.Param(ctrIdParam))
			if err != nil {
				_ = gc.Error(err)
				return
			}
			gc.String(http.StatusOK, id)
		case model.RestartingState:
			id, err := a.RestartContainer(gc.Request.Context(), gc.Param(ctrIdParam))
			if err != nil {
				_ = gc.Error(err)
				return
			}
			gc.String(http.StatusOK, id)
		default:
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(fmt.Errorf("unknown container state '%s'", ctrlReq.State))
			return
		}
	}
}
