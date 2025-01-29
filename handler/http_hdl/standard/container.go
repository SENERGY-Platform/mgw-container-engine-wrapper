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

package standard

import (
	"fmt"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/http_hdl/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

type containersQuery struct {
	Name   string `form:"name"`
	State  string `form:"state"`
	Labels string `form:"labels"`
}

type deleteContainerQuery struct {
	Force bool `form:"force"`
}

func getContainersH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, model.ContainersPath, func(c *gin.Context) {
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
		filter.Labels = util.GenLabels(util.ParseStringSlice(query.Labels, ","))
		containers, err := a.GetContainers(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, containers)
	}
}

func postContainerH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, model.ContainersPath, func(c *gin.Context) {
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

func deleteContainerH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(model.ContainersPath, ":id"), func(c *gin.Context) {
		query := deleteContainerQuery{}
		if err := c.ShouldBindQuery(&query); err != nil {
			_ = c.Error(model.NewInvalidInputError(err))
			return
		}
		if err := a.RemoveContainer(c.Request.Context(), c.Param("id"), query.Force); err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func getContainerH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(model.ContainersPath, ":id"), func(c *gin.Context) {
		container, err := a.GetContainer(c.Request.Context(), c.Param("id"))
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, container)
	}
}

func patchContainerStartH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPatch, path.Join(model.ContainersPath, ":id", model.ContainerStartPath), func(gc *gin.Context) {
		err := a.StartContainer(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func patchContainerStopH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPatch, path.Join(model.ContainersPath, ":id", model.ContainerStopPath), func(gc *gin.Context) {
		jID, err := a.StopContainer(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, jID)
	}
}

func patchContainerRestartH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPatch, path.Join(model.ContainersPath, ":id", model.ContainerRestartPath), func(gc *gin.Context) {
		jID, err := a.RestartContainer(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, jID)
	}
}

func patchContainerExecH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPatch, path.Join(model.ContainersPath, ":id", model.ContainerExecPath), func(gc *gin.Context) {
		eConf := model.ExecConfig{}
		if err := gc.ShouldBindJSON(&eConf); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		jID, err := a.ContainerExec(gc.Request.Context(), gc.Param("id"), eConf)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, jID)
	}
}
