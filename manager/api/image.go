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
	"deployment-manager/manager/api/util"
	"fmt"
	"github.com/SENERGY-Platform/mgw-deployment-manager-lib/dm-lib"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func (a *Api) GetImages(gc *gin.Context) {
	query := util.ImagesQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	filter := dm_lib.ImageFilter{Labels: util.GenLabels(query.Label)}
	images, err := a.ceHandler.ListImages(gc.Request.Context(), filter)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &images)
}

func (a *Api) PostImage(gc *gin.Context) {
	image := dm_lib.Image{}
	if err := gc.ShouldBindJSON(&image); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	if len(image.Tags) == 0 {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(fmt.Errorf("missing image reference"))
		return
	}
	if err := a.ceHandler.ImagePull(gc.Request.Context(), image.Tags[0]); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}

func (a *Api) GetImage(gc *gin.Context) {
	param, err := url.QueryUnescape(gc.Param(util.ImageParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	image, err := a.ceHandler.ImageInfo(gc.Request.Context(), param)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &image)
}

func (a *Api) DeleteImage(gc *gin.Context) {
	param, err := url.QueryUnescape(gc.Param(util.ImageParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	if err := a.ceHandler.ImageRemove(gc.Request.Context(), param); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}
