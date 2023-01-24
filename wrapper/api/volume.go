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
	"container-engine-wrapper/wrapper/api/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/wrapper/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) GetVolumes(gc *gin.Context) {
	query := util.VolumesQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	volumes, err := a.ceHandler.ListVolumes(gc.Request.Context(), model.VolumeFilter{Labels: util.GenLabels(query.Label)})
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &volumes)
}

func (a *Api) PostVolume(gc *gin.Context) {
	volume := model.Volume{}
	if err := gc.ShouldBindJSON(&volume); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	if err := a.ceHandler.VolumeCreate(gc.Request.Context(), volume); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}

func (a *Api) GetVolume(gc *gin.Context) {
	volume, err := a.ceHandler.VolumeInfo(gc.Request.Context(), gc.Param(util.VolumeParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &volume)
}

func (a *Api) DeleteVolume(gc *gin.Context) {
	if err := a.ceHandler.VolumeRemove(gc.Request.Context(), gc.Param(util.VolumeParam)); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}
