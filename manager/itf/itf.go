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
	"context"
	"github.com/SENERGY-Platform/mgw-container-engine-manager-lib/cem-lib"
	"io"
)

type ContainerEngineHandler interface {
	ListNetworks(ctx context.Context) ([]cem_lib.Network, error)
	ListContainers(ctx context.Context, filter cem_lib.ContainerFilter) ([]cem_lib.Container, error)
	ListImages(ctx context.Context, filter cem_lib.ImageFilter) ([]cem_lib.Image, error)
	ListVolumes(ctx context.Context, filter cem_lib.VolumeFilter) ([]cem_lib.Volume, error)
	NetworkInfo(ctx context.Context, id string) (cem_lib.Network, error)
	NetworkCreate(ctx context.Context, net cem_lib.Network) error
	NetworkRemove(ctx context.Context, id string) error
	ContainerInfo(ctx context.Context, id string) (cem_lib.Container, error)
	ContainerCreate(ctx context.Context, container cem_lib.Container) (id string, err error)
	ContainerRemove(ctx context.Context, id string) error
	ContainerStart(ctx context.Context, id string) error
	ContainerStop(ctx context.Context, id string) error
	ContainerRestart(ctx context.Context, id string) error
	ContainerLog(ctx context.Context, id string, logOptions cem_lib.LogOptions) (io.ReadCloser, error)
	ImageInfo(ctx context.Context, id string) (cem_lib.Image, error)
	ImagePull(ctx context.Context, id string) error
	ImageRemove(ctx context.Context, id string) error
	VolumeInfo(ctx context.Context, id string) (cem_lib.Volume, error)
	VolumeCreate(ctx context.Context, vol cem_lib.Volume) error
	VolumeRemove(ctx context.Context, id string) error
}
