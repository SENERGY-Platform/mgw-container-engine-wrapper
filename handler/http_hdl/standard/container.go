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

// getContainersH godoc
// @Summary Get containers
// @Description List all containers.
// @Tags Containers
// @Produce	json
// @Param name query string false "filter by name"
// @Param state query string false "filter by state"
// @Param labels query string false "filter by label (e.g.: l1=v1,l2=v2,l3)"
// @Success	200 {array} model.Container "containers"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers [get]
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

// postContainerH godoc
// @Summary Create container
// @Description Create a new container.
// @Tags Containers
// @Accept json
// @Produce	plain
// @Param data body model.Container true "container data"
// @Success	200 {string} string "container ID"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers [post]
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

// deleteContainerH godoc
// @Summary Delete container
// @Description Remove a container
// @Tags Containers
// @Param id path string true "container ID"
// @Param force query string false "force remove"
// @Success	200
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers/{id} [delete]
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

// getContainerH godoc
// @Summary Get container
// @Description Get a container.
// @Tags Containers
// @Produce	json
// @Param id path string true "container ID"
// @Success	200 {object} model.Container "container data"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers/{id} [get]
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

// patchContainerStartH godoc
// @Summary Start container
// @Description Start a container.
// @Tags Containers
// @Param id path string true "container ID"
// @Success	200
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers/{id}/start [patch]
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

// patchContainerStopH godoc
// @Summary Stop container
// @Description Stop a container.
// @Tags Containers
// @Produce	plain
// @Param id path string true "container ID"
// @Success	200 {string} string "job ID"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers/{id}/stop [patch]
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

// patchContainerRestartH godoc
// @Summary Restart container
// @Description Restart a container.
// @Tags Containers
// @Produce	plain
// @Param id path string true "container ID"
// @Success	200 {string} string " job ID"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers/{id}/restart [patch]
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

// patchContainerExecH godoc
// @Summary Execute command
// @Description Execute a command in a running container.
// @Tags Containers
// @Accept json
// @Produce	plain
// @Param id path string true "container ID"
// @Param cmd body model.ExecConfig true "command data"
// @Success	200 {string} string "job ID"
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /containers/{id}/exec [patch]
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
