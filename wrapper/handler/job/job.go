package job

import (
	"context"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/wrapper/model"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Job struct {
	mu    sync.RWMutex
	meta  model.Job
	tFunc func()
	ctx   context.Context
	cFunc context.CancelFunc
}

func NewJob(ctx context.Context, id uuid.UUID) *Job {
	c, cf := context.WithCancel(ctx)
	return &Job{
		meta: model.Job{
			ID:      id,
			Created: model.Time(time.Now().UTC()),
		},
		ctx:   c,
		cFunc: cf,
	}
}

func (j *Job) SetTarget(f func()) {
	j.tFunc = f
}

func (j *Job) CallTarget(cbk func()) {
	j.tFunc()
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

func (j *Job) SetStarted() {
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Started = &t
	j.mu.Unlock()
}

func (j *Job) SetCompleted() {
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Completed = &t
	j.mu.Unlock()
}

func (j *Job) SetResult(ref string, err error) {
	j.mu.Lock()
	j.meta.Ref = ref
	j.meta.Error = err.Error()
	j.mu.Unlock()
}
