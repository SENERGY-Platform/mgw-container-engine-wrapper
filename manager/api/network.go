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
	"github.com/SENERGY-Platform/mgw-container-engine-manager/manager/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) GetNetworks(gc *gin.Context) {
	networks, err := a.ceHandler.ListNetworks(gc.Request.Context())
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &networks)
}

func (a *Api) PostNetwork(gc *gin.Context) {
	network := model.Network{}
	if err := gc.ShouldBindJSON(&network); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	if err := a.ceHandler.NetworkCreate(gc.Request.Context(), network); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}

func (a *Api) GetNetwork(gc *gin.Context) {
	network, err := a.ceHandler.NetworkInfo(gc.Request.Context(), gc.Param(util.NetworkParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &network)
}

func (a *Api) DeleteNetwork(gc *gin.Context) {
	if err := a.ceHandler.NetworkRemove(gc.Request.Context(), gc.Param(util.NetworkParam)); err != nil {
		_ = gc.Error(err)
		return
	}
	gc.Status(http.StatusOK)
}
