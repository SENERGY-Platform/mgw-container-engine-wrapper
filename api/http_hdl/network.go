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

const netIdParam = "n"

func getNetworksH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		networks, err := a.GetNetworks(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, networks)
	}
}

func postNetworkH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		network := model.Network{}
		if err := gc.ShouldBindJSON(&network); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		id, err := a.CreateNetwork(gc.Request.Context(), network)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, id)
	}
}

func getNetworkH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		network, err := a.GetNetwork(gc.Request.Context(), gc.Param(netIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, network)
	}
}

func deleteNetworkH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		if err := a.RemoveNetwork(gc.Request.Context(), gc.Param(netIdParam)); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
