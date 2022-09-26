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
	"deployment-manager/manager/itf"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a Api) GetContainers(gc *gin.Context) {
	query := ContainersQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	filter := itf.ContainerFilter{Name: query.Name}
	if query.State != "" {
		cs, ok := itf.ContainerStateMap[query.State]
		if !ok {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(fmt.Errorf("unknown container state '%s'", query.State))
			return
		}
		filter.State = cs
	}
	filter.Labels = GenLabels(query.Label)
	containers, err := a.ceHandler.ListContainers(gc.Request.Context(), filter)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, &containers)
}

func (a Api) PostContainer(gc *gin.Context) {

}

func (a Api) PutContainer(gc *gin.Context) {

}

func (a Api) DeleteContainer(gc *gin.Context) {

}

func (a Api) GetContainer(gc *gin.Context) {

}

func (a Api) GetContainerLog(gc *gin.Context) {

}

func (a Api) GetImages(gc *gin.Context) {

}

func (a Api) PostImage(gc *gin.Context) {

}

func (a Api) GetImage(gc *gin.Context) {

}

func (a Api) DeleteImage(gc *gin.Context) {

}

func (a Api) GetNetworks(gc *gin.Context) {

}

func (a Api) PostNetwork(gc *gin.Context) {

}

func (a Api) GetNetwork(gc *gin.Context) {

}

func (a Api) DeleteNetwork(gc *gin.Context) {

}
