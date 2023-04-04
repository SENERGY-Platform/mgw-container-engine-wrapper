package job

import (
	"container-engine-wrapper/model"
	"context"
	"sync"
	"time"
)

type job struct {
	mu    sync.RWMutex
	meta  model.Job
	tFunc func(context.Context, context.CancelFunc) error
	ctx   context.Context
	cFunc context.CancelFunc
}

func (j *job) CallTarget(cbk func()) {
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Started = &t
	j.mu.Unlock()
	err := j.tFunc(j.ctx, j.cFunc)
	j.mu.Lock()
	if err != nil {
		j.meta.Error = err.Error()
	}
	t = model.Time(time.Now().UTC())
	j.meta.Completed = &t
	j.mu.Unlock()
	cbk()
}

func (j *job) IsCanceled() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.ctx.Err() == context.Canceled
}

func (j *job) Cancel() {
	j.cFunc()
	j.mu.Lock()
	t := model.Time(time.Now().UTC())
	j.meta.Canceled = &t
	j.mu.Unlock()
}

func (j *job) Meta() model.Job {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.meta
}
