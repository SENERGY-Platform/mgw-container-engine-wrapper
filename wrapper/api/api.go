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
	"container-engine-wrapper/wrapper/api/util"
	"container-engine-wrapper/wrapper/itf"
	"github.com/gin-gonic/gin"
)

type Api struct {
	ceHandler  itf.ContainerEngineHandler
	jobHandler itf.JobHandler
	rHeaders   RequestHeaders
}

type RequestHeaders struct {
	RequestUri string
	Uri        string
}

func New(ceHandler itf.ContainerEngineHandler, jobHandler itf.JobHandler, rHeaders RequestHeaders) *Api {
	return &Api{
		ceHandler:  ceHandler,
		jobHandler: jobHandler,
		rHeaders:   rHeaders,
	}
}

func (a *Api) SetRoutes(e *gin.Engine) {
	e.GET("/containers", a.GetContainers)
	e.POST("/containers", a.PostContainer)
	e.DELETE("/containers/:"+util.ContainerParam, a.DeleteContainer)
	e.GET("/containers/:"+util.ContainerParam, a.GetContainer)
	e.POST("/containers/:"+util.ContainerParam+"/start", a.PostContainerStart)
	e.POST("/containers/:"+util.ContainerParam+"/stop", a.PostContainerStop)
	e.POST("/containers/:"+util.ContainerParam+"/restart", a.PostContainerRestart)
	e.GET("/logs/:"+util.ContainerParam, a.GetContainerLog)
	e.GET("/images", a.GetImages)
	e.POST("/images", a.PostImage)
	e.GET("/images/:"+util.ImageParam, a.GetImage)
	e.DELETE("/images/:"+util.ImageParam, a.DeleteImage)
	e.GET("/networks", a.GetNetworks)
	e.POST("/networks", a.PostNetwork)
	e.GET("/networks/:"+util.NetworkParam, a.GetNetwork)
	e.DELETE("/networks/:"+util.NetworkParam, a.DeleteNetwork)
	e.GET("/volumes", a.GetVolumes)
	e.POST("/volumes", a.PostVolume)
	e.GET("/volumes/:"+util.VolumeParam, a.GetVolume)
	e.DELETE("/volumes/:"+util.VolumeParam, a.DeleteVolume)
	e.GET("/jobs", a.GetJobs)
	e.GET("/jobs/:"+util.JobParam, a.GetJob)
	e.POST("/jobs/:"+util.JobParam+"/cancel", a.PostJobCancel)
}
