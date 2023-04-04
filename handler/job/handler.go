package job

import (
	"container-engine-wrapper/model"
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/go-cc-job-handler/ccjh"
	"sort"
	"sync"
	"time"
)

type Handler struct {
	mu        sync.RWMutex
	ctx       context.Context
	ccHandler *ccjh.Handler
	jobs      map[string]*model.JobInternal
}

func New(ctx context.Context, ccHandler *ccjh.Handler) *Handler {
	return &Handler{
		ctx:       ctx,
		ccHandler: ccHandler,
		jobs:      make(map[string]*model.JobInternal),
	}
}

func (h *Handler) Add(id string, job *model.JobInternal) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.jobs[id]; ok {
		return errors.New("duplicate job id")
	}
	err := h.ccHandler.Add(job)
	if err != nil {
		return err
	}
	h.jobs[id] = job
	return nil
}

func (h *Handler) Get(id string) (*model.JobInternal, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	j, ok := h.jobs[id]
	if !ok {
		return nil, fmt.Errorf("%s not found", id)
	}
	return j, nil
}

func (h *Handler) List(filter model.JobOptions) []model.Job {
	var jobs []model.Job
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, v := range h.jobs {
		if check(filter, v.Meta()) {
			jobs = append(jobs, v.Meta())
		}
	}
	if filter.SortDesc {
		sort.Slice(jobs, func(i, j int) bool {
			return time.Time(jobs[i].Created).After(time.Time(jobs[j].Created))
		})
	} else {
		sort.Slice(jobs, func(i, j int) bool {
			return time.Time(jobs[i].Created).Before(time.Time(jobs[j].Created))
		})
	}
	return jobs
}

func (h *Handler) Context() context.Context {
	return h.ctx
}

func (h *Handler) PurgeJobs(maxAge int64) int {
	var l []string
	tNow := time.Now().UTC()
	h.mu.RLock()
	for k, v := range h.jobs {
		m := v.Meta()
		if v.IsCanceled() || m.Completed != nil || m.Canceled != nil {
			if tNow.Sub(time.Time(m.Created)).Microseconds() >= maxAge {
				l = append(l, k)
			}
		}
	}
	h.mu.RUnlock()
	h.mu.Lock()
	for _, id := range l {
		delete(h.jobs, id)
	}
	h.mu.Unlock()
	return len(l)
}

func check(filter model.JobOptions, job model.Job) bool {
	jC := time.Time(job.Created)
	tS := filter.Since
	tU := filter.Until
	if !tS.IsZero() && !jC.After(tS) {
		return false
	}
	if !tU.IsZero() && !jC.Before(tU) {
		return false
	}
	switch filter.Status {
	case model.JobPending:
		if job.Started != nil || job.Canceled != nil || job.Completed != nil {
			return false
		}
	case model.JobRunning:
		if job.Started == nil || job.Canceled != nil || job.Completed != nil {
			return false
		}
	case model.JobCanceled:
		if job.Canceled == nil {
			return false
		}
	case model.JobCompleted:
		if job.Completed == nil {
			return false
		}
	case model.JobError:
		if job.Completed != nil && job.Error == nil {
			return false
		}
	case model.JobOK:
		if job.Completed != nil && job.Error != nil {
			return false
		}
	}
	return true
}
