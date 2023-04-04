package job

import (
	"container-engine-wrapper/model"
	"context"
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
	jobs      map[string]*job
}

func New(ctx context.Context, ccHandler *ccjh.Handler) *Handler {
	return &Handler{
		ctx:       ctx,
		ccHandler: ccHandler,
		jobs:      make(map[string]*job),
	}
}

func (h *Handler) Create(desc string, tFunc func(context.Context, context.CancelFunc) error) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	id := uid.String()
	ctx, cf := context.WithCancel(h.ctx)
	j := job{
		meta: model.Job{
			ID:          id,
			Created:     model.Time(time.Now().UTC()),
			Description: desc,
		},
		tFunc: tFunc,
		ctx:   ctx,
		cFunc: cf,
	}
	h.mu.Lock()
	h.jobs[id] = &j
	defer h.mu.Unlock()
	return id, nil
}

func (h *Handler) Get(id string) (model.Job, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	j, ok := h.jobs[id]
	if !ok {
		return model.Job{}, fmt.Errorf("%s not found", id)
	}
	return j.Meta(), nil
}

func (h *Handler) Cancel(id string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	j, ok := h.jobs[id]
	if !ok {
		return fmt.Errorf("%s not found", id)
	}
	j.Cancel()
	return nil
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
