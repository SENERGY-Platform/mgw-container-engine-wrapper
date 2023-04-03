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
	"container-engine-wrapper/api/util"
	"container-engine-wrapper/itf"
	"container-engine-wrapper/model"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (a *Api) GetImages(gc *gin.Context) {
	query := util.ImagesQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	filter := itf.ImageFilter{Labels: util.GenLabels(query.Label)}
	images, err := a.ceHandler.ListImages(gc.Request.Context(), filter)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &images)
}

func (a *Api) PostImage(gc *gin.Context) {
	req := model.ImagesPostRequest{}
	if err := gc.ShouldBindJSON(&req); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	img := req.Image
	jID, err := uuid.NewRandom()
	if err != nil {
		_ = gc.Error(err)
		return
	}
	ctx, cf := context.WithCancel(a.jobHandler.Context())
	j := itf.NewJob(ctx, cf, jID.String(), model.JobOrgRequest{
		Method: gc.Request.Method,
		Uri:    gc.Request.RequestURI,
		Body:   req,
	})
	j.SetTarget(func() {
		defer cf()
		e := a.ceHandler.ImagePull(ctx, img)
		if e == nil {
			e = ctx.Err()
		}
		j.SetError(e)
	})
	err = a.jobHandler.Add(jID, j)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	rUri := gc.GetHeader(a.rHeaders.RequestUri)
	uri := gc.GetHeader(a.rHeaders.Uri)
	if rUri != "" || uri != "" {
		gc.Redirect(http.StatusSeeOther, strings.Replace(rUri, uri, "/", 1)+"jobs/"+jID.String())
	} else {
		gc.Status(http.StatusOK)
	}
}

func (a *Api) GetImage(gc *gin.Context) {
	image, err := a.ceHandler.ImageInfo(gc.Request.Context(), gc.Param(util.ImageParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &image)
}

func (a *Api) DeleteImage(gc *gin.Context) {
	if err := a.ceHandler.ImageRemove(gc.Request.Context(), gc.Param(util.ImageParam)); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}
