package job

import (
	"container-engine-wrapper/itf"
	"container-engine-wrapper/model"
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/go-cc-job-handler/ccjh"
	"github.com/google/uuid"
	"sort"
	"sync"
	"time"
)

type Handler struct {
	mu        sync.RWMutex
	ctx       context.Context
	ccHandler *ccjh.Handler
	jobs      map[uuid.UUID]*itf.Job
}

func New(ctx context.Context, ccHandler *ccjh.Handler) *Handler {
	return &Handler{
		ctx:       ctx,
		ccHandler: ccHandler,
		jobs:      make(map[uuid.UUID]*itf.Job),
	}
}

func (h *Handler) Add(id uuid.UUID, job *itf.Job) error {
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

func (h *Handler) Get(id uuid.UUID) (*itf.Job, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	j, ok := h.jobs[id]
	if !ok {
		return nil, fmt.Errorf("%s not found", id)
	}
	return j, nil
}

func (h *Handler) List(filter itf.JobOptions) []model.Job {
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
	var l []uuid.UUID
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

func check(filter itf.JobOptions, job model.Job) bool {
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
	case itf.JobPending:
		if job.Started != nil || job.Canceled != nil || job.Completed != nil {
			return false
		}
	case itf.JobRunning:
		if job.Started == nil || job.Canceled != nil || job.Completed != nil {
			return false
		}
	case itf.JobCanceled:
		if job.Canceled == nil {
			return false
		}
	case itf.JobCompleted:
		if job.Completed == nil {
			return false
		}
	case itf.JobError:
		if job.Completed != nil && job.Error == nil {
			return false
		}
	case itf.JobOK:
		if job.Completed != nil && job.Error != nil {
			return false
		}
	}
	return true
}
