/*
 * Copyright 2022 InfAI (CC SES)
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

package itf

import (
	"container-engine-wrapper/model"
	"context"
	"github.com/google/uuid"
	"io"
)

type ContainerEngineHandler interface {
	ListNetworks(ctx context.Context) ([]model.Network, error)
	ListContainers(ctx context.Context, filter ContainerFilter) ([]model.Container, error)
	ListImages(ctx context.Context, filter ImageFilter) ([]model.Image, error)
	ListVolumes(ctx context.Context, filter VolumeFilter) ([]model.Volume, error)
	NetworkInfo(ctx context.Context, id string) (model.Network, error)
	NetworkCreate(ctx context.Context, net model.Network) error
	NetworkRemove(ctx context.Context, id string) error
	ContainerInfo(ctx context.Context, id string) (model.Container, error)
	ContainerCreate(ctx context.Context, container model.Container) (id string, err error)
	ContainerRemove(ctx context.Context, id string) error
	ContainerStart(ctx context.Context, id string) error
	ContainerStop(ctx context.Context, id string) error
	ContainerRestart(ctx context.Context, id string) error
	ContainerLog(ctx context.Context, id string, logOptions LogOptions) (io.ReadCloser, error)
	ImageInfo(ctx context.Context, id string) (model.Image, error)
	ImagePull(ctx context.Context, id string) error
	ImageRemove(ctx context.Context, id string) error
	VolumeInfo(ctx context.Context, id string) (model.Volume, error)
	VolumeCreate(ctx context.Context, vol model.Volume) error
	VolumeRemove(ctx context.Context, id string) error
}

type JobHandler interface {
	Add(id uuid.UUID, job *Job) error
	Get(id uuid.UUID) (*Job, error)
	List(filter JobOptions) []model.Job
	Context() context.Context
}
