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
	"container-engine-wrapper/itf"
	"container-engine-wrapper/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const volIdParam = "v"

type volumesQuery struct {
	Label []string `form:"label"`
}

func getVolumesH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := volumesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		volumes, err := a.GetVolumes(gc.Request.Context(), model.VolumeFilter{Labels: GenLabels(query.Label)})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, volumes)
	}
}

func postVolumeH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		volume := model.Volume{}
		if err := gc.ShouldBindJSON(&volume); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		if err := a.CreateVolume(gc.Request.Context(), volume); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func getVolumeH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		volume, err := a.GetVolume(gc.Request.Context(), gc.Param(volIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, volume)
	}
}

func deleteVolumeH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		if err := a.RemoveVolume(gc.Request.Context(), gc.Param(volIdParam)); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
