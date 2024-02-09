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
	"sort"
)

func SetRoutes(e *gin.Engine, a lib.Api) {
	standardGrp := e.Group("")
	restrictedGrp := e.Group(model.RestrictedPath)
	setSharedRoutes(a, standardGrp, restrictedGrp)
	setContainersRoutes(a, standardGrp.Group(model.ContainersPath))
	setImagesRoutes(a, standardGrp.Group(model.ImagesPath))
	setNetworksRoutes(a, standardGrp.Group(model.NetworksPath))
	setVolumesRoutes(a, standardGrp.Group(model.VolumesPath))
	setJobsRoutes(a, standardGrp.Group(model.JobsPath))
}

func GetRoutes(e *gin.Engine) [][2]string {
	routes := e.Routes()
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})
	var rInfo [][2]string
	for _, info := range routes {
		rInfo = append(rInfo, [2]string{info.Method, info.Path})
	}
	return rInfo
}

func setSharedRoutes(a lib.Api, rGroups ...*gin.RouterGroup) {
	for _, rg := range rGroups {
		rg.GET(model.SrvInfoPath, getSrvInfoH(a))
		rg.GET(model.ContainerLogsPath+"/:"+ctrIdParam, getContainerLogH(a))
	}
}

func setContainersRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET("", getContainersH(a))
	rg.POST("", postContainerH(a))
	rg.DELETE(":"+ctrIdParam, deleteContainerH(a))
	rg.GET(":"+ctrIdParam, getContainerH(a))
	rg.PATCH(":"+ctrIdParam+"/"+model.ContainerStartPath, patchContainerStartH(a))
	rg.PATCH(":"+ctrIdParam+"/"+model.ContainerStopPath, patchContainerStopH(a))
	rg.PATCH(":"+ctrIdParam+"/"+model.ContainerRestartPath, patchContainerRestartH(a))
	rg.PATCH(":"+ctrIdParam+"/"+model.ContainerExecPath, patchContainerExecH(a))
}

func setImagesRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET("", getImagesH(a))
	rg.POST("", postImageH(a))
	rg.GET(":"+imgIdParam, getImageH(a))
	rg.DELETE(":"+imgIdParam, deleteImageH(a))
}

func setNetworksRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET("", getNetworksH(a))
	rg.POST("", postNetworkH(a))
	rg.GET(":"+netIdParam, getNetworkH(a))
	rg.DELETE(":"+netIdParam, deleteNetworkH(a))
}

func setVolumesRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET("", getVolumesH(a))
	rg.POST("", postVolumeH(a))
	rg.GET(":"+volIdParam, getVolumeH(a))
	rg.DELETE(":"+volIdParam, deleteVolumeH(a))
}

func setJobsRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET("", getJobsH(a))
	rg.GET(":"+jobIdParam, getJobH(a))
	rg.PATCH(":"+jobIdParam+"/"+model.JobsCancelPath, patchJobCancelH(a))
}
