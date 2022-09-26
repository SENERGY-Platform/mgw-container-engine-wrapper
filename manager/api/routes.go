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
	"github.com/gin-gonic/gin"
)

func SetRoutes(e *gin.Engine, a *Api) {
	e.GET("/containers", a.GetContainers)
	e.POST("/containers", a.PostContainer)
	e.PUT("/containers/:container", a.PutContainer)
	e.DELETE("/containers/:container", a.DeleteContainer)
	e.GET("/containers/:container", a.GetContainer)
	e.GET("/containers/:container/log", a.GetContainerLog)
	e.GET("/images", a.GetImages)
	e.POST("/images", a.PostImage)
	e.GET("/images/:image", a.GetImage)
	e.DELETE("/images/:image", a.DeleteImage)
	e.GET("/networks", a.GetNetworks)
	e.POST("/networks", a.PostNetwork)
	e.GET("/networks/:network", a.GetNetwork)
	e.DELETE("/networks/:network", a.DeleteNetwork)
}
