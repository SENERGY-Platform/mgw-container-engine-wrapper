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

package wrapper

import (
	"github.com/SENERGY-Platform/mgw-go-service-base/job-hdl"
	"github.com/SENERGY-Platform/mgw-go-service-base/srv-info-hdl"
)

type Wrapper struct {
	ceHandler  ContainerEngineHandler
	jobHandler job_hdl.JobHandler
	srvInfoHdl srv_info_hdl.SrvInfoHandler
}

func New(ceHandler ContainerEngineHandler, jobHandler job_hdl.JobHandler, srvInfoHandler srv_info_hdl.SrvInfoHandler) *Wrapper {
	return &Wrapper{
		ceHandler:  ceHandler,
		jobHandler: jobHandler,
		srvInfoHdl: srvInfoHandler,
	}
}
