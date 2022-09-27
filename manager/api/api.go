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
	"github.com/gin-gonic/gin"
)

type Api struct {
	ceHandler itf.ContainerEngineHandler
}

func New(ceHandler itf.ContainerEngineHandler) *Api {
	return &Api{
		ceHandler: ceHandler,
	}
}

func (a Api) SetRoutes(e *gin.Engine) {
	e.GET("/containers", a.GetContainers)
	e.POST("/containers", a.PostContainer)
	e.PUT("/containers/:"+containerParam, a.PutContainer)
	e.DELETE("/containers/:"+containerParam, a.DeleteContainer)
	e.GET("/containers/:"+containerParam, a.GetContainer)
	e.GET("/containers/:"+containerParam+"/log", a.GetContainerLog)
	e.GET("/images", a.GetImages)
	e.POST("/images", a.PostImage)
	e.GET("/images/:"+imageParam, a.GetImage)
	e.DELETE("/images/:"+imageParam, a.DeleteImage)
	e.GET("/networks", a.GetNetworks)
	e.POST("/networks", a.PostNetwork)
	e.GET("/networks/:"+networkParam, a.GetNetwork)
	e.DELETE("/networks/:"+networkParam, a.DeleteNetwork)
}
