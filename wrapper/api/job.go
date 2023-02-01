package api

import (
	"container-engine-wrapper/wrapper/api/util"
	"container-engine-wrapper/wrapper/itf"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *Api) GetJobs(gc *gin.Context) {
	query := util.JobsQuery{}
	if err := gc.ShouldBindQuery(&query); err != nil {
		gc.Status(http.StatusBadRequest)
		_ = gc.Error(err)
		return
	}
	jobOptions := itf.JobOptions{SortDesc: query.SortDesc}
	if query.State != "" {
		_, ok := itf.JobStateMap[query.State]
		if !ok {
			gc.Status(http.StatusBadRequest)
			_ = gc.Error(fmt.Errorf("unknown job state '%s'", query.State))
			return
		}
		jobOptions.State = query.State
	}
	if query.Since > 0 {
		jobOptions.Since = time.UnixMicro(query.Since)
	}
	if query.Until > 0 {
		jobOptions.Until = time.UnixMicro(query.Until)
	}
	jobs := a.jobHandler.List(jobOptions)
	gc.JSON(http.StatusOK, jobs)
}

func (a *Api) GetJob(gc *gin.Context) {
	id, err := uuid.Parse(gc.Param(util.JobParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	j, err := a.jobHandler.Get(id)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	gc.JSON(http.StatusOK, j.Meta())
}

func (a *Api) PostJobCancel(gc *gin.Context) {
	id, err := uuid.Parse(gc.Param(util.JobParam))
	if err != nil {
		_ = gc.Error(err)
		return
	}
	j, err := a.jobHandler.Get(id)
	if err != nil {
		_ = gc.Error(err)
		return
	}
	j.Cancel()
	gc.Status(http.StatusOK)
}
