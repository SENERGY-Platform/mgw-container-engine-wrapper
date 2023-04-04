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

func getJobsH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := JobsQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(err)
			return
		}
		jobOptions := model.JobOptions{SortDesc: query.SortDesc}
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
		j, err := a.GetJob(gc.Request.Context(), gc.Param(jobIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, j.Meta())
	}
}

func postJobCancelH(a itf.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		j, err := a.GetJob(gc.Request.Context(), gc.Param(jobIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		j.Cancel()
		gc.Status(http.StatusOK)
	}
}
