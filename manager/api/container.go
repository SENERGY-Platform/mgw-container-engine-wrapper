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

package api

import (
	"container-engine-manager/manager/api/util"
	"fmt"
	"github.com/SENERGY-Platform/mgw-container-engine-manager/manager/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func (a *Api) GetContainers(gc *gin.Context) {
	query := util.ContainersQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	filter := model.ContainerFilter{Name: query.Name}
	if query.State != "" {
		_, ok := model.ContainerStateMap[query.State]
		if !ok {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(fmt.Errorf("unknown container state '%s'", query.State))
			return
		}
		filter.State = query.State
	}
	filter.Labels = util.GenLabels(query.Label)
	containers, err := a.ceHandler.ListContainers(gc.Request.Context(), filter)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &containers)
}

func (a *Api) PostContainer(gc *gin.Context) {
	container := model.Container{}
	if err := gc.ShouldBindJSON(&container); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	id, err := a.ceHandler.ContainerCreate(gc.Request.Context(), container)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &model.ContainersPostResponse{ID: id})
}

func (a *Api) PostContainerCtrl(gc *gin.Context) {
	ctrl := model.ContainerCtrlPostRequest{}
	if err := gc.ShouldBindJSON(&ctrl); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	switch ctrl.State {
	case model.ContainerStart:
		if err := a.ceHandler.ContainerStart(gc.Request.Context(), gc.Param(util.ContainerParam)); err != nil {
			_ = gc.Error(err)
			return
		}
	case model.ContainerRestart:
		if err := a.ceHandler.ContainerRestart(gc.Request.Context(), gc.Param(util.ContainerParam)); err != nil {
			_ = gc.Error(err)
			return
		}
	case model.ContainerStop:
		if err := a.ceHandler.ContainerStop(gc.Request.Context(), gc.Param(util.ContainerParam)); err != nil {
			_ = gc.Error(err)
			return
		}
	default:
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(fmt.Errorf("invalid container state '%s'", ctrl.State))
	}
	gc.Status(http.StatusOK)
}

func (a *Api) DeleteContainer(gc *gin.Context) {
	if err := a.ceHandler.ContainerRemove(gc.Request.Context(), gc.Param(util.ContainerParam)); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}

func (a *Api) GetContainer(gc *gin.Context) {
	container, err := a.ceHandler.ContainerInfo(gc.Request.Context(), gc.Param(util.ContainerParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &container)
}

func (a *Api) GetContainerLog(gc *gin.Context) {
	query := util.ContainerLogQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	logOptions := model.LogOptions{MaxLines: query.MaxLines}
	if query.Since != "" {
		since, err := time.Parse(time.RFC3339Nano, query.Since)
		if err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		logOptions.Since = &since
	}
	if query.Until != "" {
		until, err := time.Parse(time.RFC3339Nano, query.Until)
		if err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		logOptions.Until = &until
	}
	rc, err := a.ceHandler.ContainerLog(gc.Request.Context(), gc.Param(util.ContainerParam), logOptions)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	defer rc.Close()
	gc.Status(http.StatusOK)
	gc.Header("Transfer-Encoding", "chunked")
	gc.Header("Content-Type", gin.MIMEPlain)
	for {
		var b = make([]byte, 204800)
		n, rErr := rc.Read(b)
		if rErr != nil {
			if rErr == io.EOF {
				if n > 0 {
					_, wErr := gc.Writer.Write(b[:n])
					if wErr != nil {
						gc.Status(http.StatusInternalServerError)
						_ = gc.Error(wErr)
						return
					}
					gc.Writer.Flush()
				}
				break
			}
			gc.Status(http.StatusInternalServerError)
			_ = gc.Error(rErr)
			return
		}
		_, wErr := gc.Writer.Write(b)
		if wErr != nil {
			gc.Status(http.StatusInternalServerError)
			_ = gc.Error(wErr)
			return
		}
		gc.Writer.Flush()
	}
}
