/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"context"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"io"
)

type Api interface {
	GetContainers(ctx context.Context, filter model.ContainerFilter) ([]model.Container, error)
	GetContainer(ctx context.Context, id string) (model.Container, error)
	CreateContainer(ctx context.Context, container model.Container) (id string, err error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) (jobId string, err error)
	RestartContainer(ctx context.Context, id string) (jobId string, err error)
	RemoveContainer(ctx context.Context, id string) error
	GetContainerLog(ctx context.Context, id string, logOptions model.LogFilter) (io.ReadCloser, error)
	GetImages(ctx context.Context, filter model.ImageFilter) ([]model.Image, error)
	GetImage(ctx context.Context, id string) (model.Image, error)
	AddImage(ctx context.Context, img string) (jobId string, err error)
	RemoveImage(ctx context.Context, id string) error
	GetNetworks(ctx context.Context) ([]model.Network, error)
	GetNetwork(ctx context.Context, id string) (model.Network, error)
	CreateNetwork(ctx context.Context, net model.Network) error
	RemoveNetwork(ctx context.Context, id string) error
	GetVolumes(ctx context.Context, filter model.VolumeFilter) ([]model.Volume, error)
	GetVolume(ctx context.Context, id string) (model.Volume, error)
	CreateVolume(ctx context.Context, vol model.Volume) error
	RemoveVolume(ctx context.Context, id string) error
	GetJobs(ctx context.Context, filter model.JobFilter) ([]model.Job, error)
	GetJob(ctx context.Context, id string) (model.Job, error)
	CancelJob(ctx context.Context, id string) error
}
