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
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

// getNetworksH godoc
// @Summary Get networks
// @Description List all container networks.
// @Tags Networks
// @Produce	json
// @Success	200 {array} model.Network "networks"
// @Failure	500 {string} string "error message"
// @Router /networks [get]
func getNetworksH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, model.NetworksPath, func(gc *gin.Context) {
		networks, err := a.GetNetworks(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, networks)
	}
}

// postNetworkH godoc
// @Summary Create network
// @Description Add a new container network.
// @Tags Networks
// @Accept json
// @Produce	plain
// @Param data body model.Network true "network data"
// @Success	200 {string} string "network ID"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /networks [post]
func postNetworkH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, model.NetworksPath, func(gc *gin.Context) {
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

// getNetworkH godoc
// @Summary Get network
// @Description Get a container network.
// @Tags Networks
// @Produce	json
// @Param id path string true "network ID"
// @Success	200 {object} model.Network "network info"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /networks/{id} [get]
func getNetworkH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(model.NetworksPath, ":id"), func(gc *gin.Context) {
		network, err := a.GetNetwork(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, network)
	}
}

// deleteNetworkH godoc
// @Summary Delete network
// @Description Remove a container network.
// @Tags Networks
// @Param id path string true "network ID"
// @Success	200
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /networks/{id} [delete]
func deleteNetworkH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(model.NetworksPath, ":id"), func(gc *gin.Context) {
		if err := a.RemoveNetwork(gc.Request.Context(), gc.Param("id")); err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
