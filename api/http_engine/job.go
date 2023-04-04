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
	"container-engine-wrapper/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const jobIdParam = "j"

type jobsQuery struct {
	Status   string `form:"status"`
	SortDesc bool   `form:"sort_desc"`
	Since    int64  `form:"since"`
	Until    int64  `form:"until"`
}

func getJobsH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := jobsQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		jobOptions := model.JobFilter{SortDesc: query.SortDesc}
		if query.Status != "" {
			_, ok := model.JobStateMap[query.Status]
			if !ok {
				gc.Status(http.StatusBadRequest)
				_ = gc.Error(fmt.Errorf("unknown job state '%s'", query.Status))
				return
			}
			jobOptions.Status = query.Status
		}
		if query.Since > 0 {
			jobOptions.Since = time.UnixMicro(query.Since)
		}
		if query.Until > 0 {
			jobOptions.Until = time.UnixMicro(query.Until)
		}
		jobs := a.GetJobs(gc.Request.Context(), jobOptions)
		gc.JSON(http.StatusOK, jobs)
	}
}

func getJobH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		job, err := a.GetJob(gc.Request.Context(), gc.Param(jobIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, job)
	}
}

func postJobCancelH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		err := a.CancelJob(gc.Request.Context(), gc.Param(jobIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
