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
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a Api) GetNetworks(gc *gin.Context) {
	networks, err := a.ceHandler.ListNetworks(gc.Request.Context())
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &networks)
}

func (a Api) PostNetwork(gc *gin.Context) {

}

func (a Api) GetNetwork(gc *gin.Context) {

}

func (a Api) DeleteNetwork(gc *gin.Context) {

}
