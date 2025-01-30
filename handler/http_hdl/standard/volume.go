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
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/http_hdl/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

type volumesQuery struct {
	Labels string `form:"labels"`
}

type deleteVolumeQuery struct {
	Force bool `form:"force"`
}

// getVolumesH godoc
// @Summary Get volumes
// @Description List all storage volumes.
// @Tags Volumes
// @Produce	json
// @Param labels query string false "filter by label (e.g.: l1=v1,l2=v2,l3)"
// @Success	200 {array} model.Volume "volumes"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /volumes [get]
func getVolumesH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, model.VolumesPath, func(gc *gin.Context) {
		query := volumesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		volumes, err := a.GetVolumes(gc.Request.Context(), model.VolumeFilter{Labels: util.GenLabels(util.ParseStringSlice(query.Labels, ","))})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, volumes)
	}
}

// postVolumeH godoc
// @Summary Create volume
// @Description Create a new storage volume.
// @Tags Volumes
// @Accept json
// @Produce	plain
// @Param data body model.Volume true "volume data"
// @Success	200 {string} string "volume ID"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /volumes [post]
func postVolumeH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, model.VolumesPath, func(gc *gin.Context) {
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

// getVolumeH godoc
// @Summary Get volume
// @Description Get storage volume info.
// @Tags Volumes
// @Produce	json
// @Param id path string true "volume ID"
// @Success	200 {object} model.Volume "volume data"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /volumes/{id} [get]
func getVolumeH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(model.VolumesPath, ":id"), func(gc *gin.Context) {
		volume, err := a.GetVolume(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, volume)
	}
}

// deleteVolumeH godoc
// @Summary Delete volume
// @Description Remove a storage volume.
// @Tags Volumes
// @Param id path string true "volume ID"
// @Param force query string false "force delete"
// @Success	200
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /volumes/{id} [delete]
func deleteVolumeH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(model.VolumesPath, ":id"), func(gc *gin.Context) {
		query := deleteVolumeQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		if err := a.RemoveVolume(gc.Request.Context(), gc.Param("id"), query.Force); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
