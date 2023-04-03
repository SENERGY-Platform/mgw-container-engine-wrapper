package itf

import (
	"container-engine-wrapper/wrapper/model"
	"context"
	"sync"
	"time"
)

type ContainerFilter struct {
	Name   string
	State  model.ContainerState
	Labels map[string]string
}

type VolumeFilter struct {
	Labels map[string]string
}

type ImageFilter struct {
	Labels map[string]string
}

type LogOptions struct {
	MaxLines int
	Since    time.Time
	Until    time.Time
}

type Job struct {
	mu    sync.RWMutex
	meta  model.Job
	tFunc func()
	ctx   context.Context
	cFunc context.CancelFunc
}

type JobStatus = string

type SortDirection = string

type JobOptions struct {
	Status   JobStatus
	SortDesc bool
	Since    time.Time
	Until    time.Time
}
