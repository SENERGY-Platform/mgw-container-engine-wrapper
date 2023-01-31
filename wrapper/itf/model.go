package itf

import "github.com/SENERGY-Platform/mgw-container-engine-wrapper/model"

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
	Since    *model.Time
	Until    *model.Time
}
