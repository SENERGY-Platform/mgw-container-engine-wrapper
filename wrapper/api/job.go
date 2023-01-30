package api

import (
	"container-engine-wrapper/wrapper/api/util"
	"container-engine-wrapper/wrapper/handler/job"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/wrapper/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (a *Api) GetJobs(gc *gin.Context) {
	var jobs []model.Job
	a.jobHandler.Range(func(_ uuid.UUID, v *job.Job) bool {
		jobs = append(jobs, v.Meta())
		return true
	})
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
	gc.JSON(http.StatusOK, &j)
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
