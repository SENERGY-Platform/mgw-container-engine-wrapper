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

package http_engine

import (
	"container-engine-wrapper/itf"
	"github.com/gin-gonic/gin"
)

func SetRoutes(e *gin.Engine, a itf.Api) {
	e.GET("/containers", getContainersH(a))
	e.POST("/containers", postContainerH(a))
	e.DELETE("/containers/:"+ctrIdParam, deleteContainerH(a))
	e.GET("/containers/:"+ctrIdParam, getContainerH(a))
	e.POST("/containers/:"+ctrIdParam+"/ctrl", postContainerCtrlH(a))
	e.GET("/logs/:"+ctrIdParam, getContainerLogH(a))
	e.GET("/images", getImagesH(a))
	e.POST("/images", postImageH(a))
	e.GET("/images/:"+imgIdParam, getImageH(a))
	e.DELETE("/images/:"+imgIdParam, deleteImageH(a))
	e.GET("/networks", getNetworksH(a))
	e.POST("/networks", postNetworkH(a))
	e.GET("/networks/:"+netIdParam, getNetworkH(a))
	e.DELETE("/networks/:"+netIdParam, deleteNetworkH(a))
	e.GET("/volumes", getVolumesH(a))
	e.POST("/volumes", postVolumeH(a))
	e.GET("/volumes/:"+volIdParam, getVolumeH(a))
	e.DELETE("/volumes/:"+volIdParam, deleteVolumeH(a))
	e.GET("/jobs", getJobsH(a))
	e.GET("/jobs/:"+jobIdParam, getJobH(a))
	e.POST("/jobs/:"+jobIdParam+"/cancel", postJobCancelH(a))
}
