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

const imgIdParam = "i"

func getImagesH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := ImagesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		filter := model.ImageFilter{Labels: GenLabels(query.Label)}
		images, err := a.GetImages(gc.Request.Context(), filter)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, images)
	}
}

func postImageH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		req := model.ImageRequest{}
		if err := gc.ShouldBindJSON(&req); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		id, err := a.AddImage(gc.Request.Context(), req.Image)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, id)
	}
}

func getImageH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		image, err := a.GetImage(gc.Request.Context(), gc.Param(imgIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, image)
	}
}

func deleteImageH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		if err := a.RemoveImage(gc.Request.Context(), gc.Param(imgIdParam)); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
