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
)

func SetRoutes(e *gin.Engine, a lib.Api) {
	e.GET("/"+model.ContainersPath, getContainersH(a))
	e.POST("/"+model.ContainersPath, postContainerH(a))
	e.DELETE("/"+model.ContainersPath+"/:"+ctrIdParam, deleteContainerH(a))
	e.GET("/"+model.ContainersPath+"/:"+ctrIdParam, getContainerH(a))
	e.POST("/"+model.ContainersPath+"/:"+ctrIdParam+"/"+model.ContainerCtrlPath, postContainerCtrlH(a))
	e.GET("/"+model.ContainerLogsPath+"/:"+ctrIdParam, getContainerLogH(a))
	e.GET("/"+model.ImagesPath, getImagesH(a))
	e.POST("/"+model.ImagesPath, postImageH(a))
	e.GET("/"+model.ImagesPath+"/:"+imgIdParam, getImageH(a))
	e.DELETE("/"+model.ImagesPath+"/:"+imgIdParam, deleteImageH(a))
	e.GET("/"+model.NetworksPath, getNetworksH(a))
	e.POST("/"+model.NetworksPath, postNetworkH(a))
	e.GET("/"+model.NetworksPath+"/:"+netIdParam, getNetworkH(a))
	e.DELETE("/"+model.NetworksPath+"/:"+netIdParam, deleteNetworkH(a))
	e.GET("/"+model.VolumesPath, getVolumesH(a))
	e.POST("/"+model.VolumesPath, postVolumeH(a))
	e.GET("/"+model.VolumesPath+"/:"+volIdParam, getVolumeH(a))
	e.DELETE("/"+model.VolumesPath+"/:"+volIdParam, deleteVolumeH(a))
	e.GET("/"+model.JobsPath, getJobsH(a))
	e.GET("/"+model.JobsPath+"/:"+jobIdParam, getJobH(a))
	e.POST("/"+model.JobsPath+"/:"+jobIdParam+"/"+model.JobsCancelPath, postJobCancelH(a))
}
