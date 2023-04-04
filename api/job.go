package api

import (
	"container-engine-wrapper/itf"
	"container-engine-wrapper/model"
	"context"
)

func (a *Api) GetJobs(_ context.Context, filter itf.JobOptions) []model.Job {
	return a.jobHandler.List(filter)
}

func (a *Api) GetJob(_ context.Context, id string) (*itf.Job, error) {
	return a.jobHandler.Get(id)
}

func (a *Api) CancelJob(_ context.Context, id string) error {
	return a.jobHandler.Cancel(id)
}
