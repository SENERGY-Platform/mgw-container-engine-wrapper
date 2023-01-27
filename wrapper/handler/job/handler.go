package job

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/go-cc-job-handler/ccjh"
	"github.com/google/uuid"
	"sync"
)

type Handler struct {
	mu        sync.RWMutex
	ctx       context.Context
	ccHandler *ccjh.Handler
	jobs      map[uuid.UUID]*Job
}

func New(ctx context.Context, ccHandler *ccjh.Handler) *Handler {
	return &Handler{
		ctx:       ctx,
		ccHandler: ccHandler,
		jobs:      make(map[uuid.UUID]*Job),
	}
}

func (h *Handler) Add(job *Job) error {
	h.mu.Lock()
	h.jobs[job.meta.ID] = job
	h.mu.Unlock()
	return h.ccHandler.Add(job)
}

func (h *Handler) Get(id uuid.UUID) (*Job, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	j, ok := h.jobs[id]
	if !ok {
		return nil, fmt.Errorf("%s not found", id)
	}
	return j, nil
}

func (h *Handler) Range(f func(k uuid.UUID, v *Job) bool) {
	h.mu.RLock()
	for k, v := range h.jobs {
		if !f(k, v) {
			break
		}
	}
	h.mu.RUnlock()
}

func (h *Handler) Context() context.Context {
	return h.ctx
}
