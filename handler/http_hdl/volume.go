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
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const volIdParam = "v"

type volumesQuery struct {
	Labels string `form:"labels"`
}

type deleteVolumeQuery struct {
	Force bool `form:"force"`
}

func getVolumesH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := volumesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		volumes, err := a.GetVolumes(gc.Request.Context(), model.VolumeFilter{Labels: genLabels(parseStringSlice(query.Labels, ","))})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, volumes)
	}
}

func postVolumeH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		volume := model.Volume{}
		if err := gc.ShouldBindJSON(&volume); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		id, err := a.CreateVolume(gc.Request.Context(), volume)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, id)
	}
}

func getVolumeH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		volume, err := a.GetVolume(gc.Request.Context(), gc.Param(volIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, volume)
	}
}

func deleteVolumeH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := deleteVolumeQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		if err := a.RemoveVolume(gc.Request.Context(), gc.Param(volIdParam), query.Force); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
