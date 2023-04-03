package itf

import (
	"container-engine-wrapper/model"
	"context"
	"time"
)

func NewJob(ctx context.Context, cf context.CancelFunc, id string, req model.JobOrgRequest) *Job {
	return &Job{
		meta: model.Job{
			ID:      id,
			Request: req,
			Created: model.Time(time.Now().UTC()),
		},
		ctx:   ctx,
		cFunc: cf,
	}
}

func (j *Job) SetTarget(f func()) {
	j.tFunc = f
}

func (j *Job) CallTarget(cbk func()) {
	j.setStarted()
	j.tFunc()
	j.setCompleted()
	cbk()
}

func (j *Job) IsCanceled() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.ctx.Err() == context.Canceled
}

func (j *Job) Cancel() {
	j.cFunc()
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Canceled = &t
	j.mu.Unlock()
}

func (j *Job) Meta() model.Job {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.meta
}

func (j *Job) setStarted() {
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Started = &t
	j.mu.Unlock()
}

func (j *Job) setCompleted() {
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Completed = &t
	j.mu.Unlock()
}

func (j *Job) SetError(err error) {
	j.mu.Lock()
	if err != nil {
		j.meta.Error = err.Error()
	}
	j.mu.Unlock()
}
