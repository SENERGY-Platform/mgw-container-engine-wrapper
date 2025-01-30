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

type imagesQuery struct {
	Name   string `form:"name"`
	Tag    string `form:"tag"`
	Labels string `form:"labels"`
}

// getImagesH godoc
// @Summary Get images
// @Description List all container images.
// @Tags Images
// @Produce	json
// @Param name query string false "filter by name"
// @Param tag query string false "filter by image tag"
// @Param labels query string false "filter by labels (e.g. l1=v1,l2=v2,l3)"
// @Success	200 {array} model.Image "images"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /images [get]
func getImagesH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, model.ImagesPath, func(gc *gin.Context) {
		query := imagesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		filter := model.ImageFilter{
			Name:   query.Name,
			Tag:    query.Tag,
			Labels: util.GenLabels(util.ParseStringSlice(query.Labels, ",")),
		}
		images, err := a.GetImages(gc.Request.Context(), filter)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, images)
	}
}

// postImageH godoc
// @Summary Add image
// @Description Download a container image.
// @Tags Images
// @Accept json
// @Produce	json,plain
// @Param data body model.ImageRequest true "image data"
// @Success	200 {string} string "job ID"
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /images [post]
func postImageH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, model.ImagesPath, func(gc *gin.Context) {
		req := model.ImageRequest{}
		if err := gc.ShouldBindJSON(&req); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		jID, err := a.AddImage(gc.Request.Context(), req.Image)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, jID)
	}
}

// getImageH godoc
// @Summary Get image
// @Description Get container image info.
// @Tags Images
// @Produce	json
// @Param id path string true "image ID"
// @Success	200 {object} model.Image "image data"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /images/{id} [get]
func getImageH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(model.ImagesPath, ":id"), func(gc *gin.Context) {
		image, err := a.GetImage(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, image)
	}
}

// deleteImageH godoc
// @Summary Delete image
// @Description Remove a container image.
// @Tags Images
// @Param id path string true "image ID"
// @Success	200
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /images/{id} [delete]
func deleteImageH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(model.ImagesPath, ":id"), func(gc *gin.Context) {
		if err := a.RemoveImage(gc.Request.Context(), gc.Param("id")); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
