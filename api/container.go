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

package api

import (
	"container-engine-wrapper/model"
	"context"
	"fmt"
	"io"
)

func (a *Api) GetContainers(ctx context.Context, filter model.ContainerFilter) ([]model.Container, error) {
	return a.ceHandler.ListContainers(ctx, filter)
}

func (a *Api) GetContainer(ctx context.Context, id string) (model.Container, error) {
	return a.ceHandler.ContainerInfo(ctx, id)
}

func (a *Api) CreateContainer(ctx context.Context, container model.Container) (string, error) {
	return a.ceHandler.ContainerCreate(ctx, container)
}

func (a *Api) StartContainer(ctx context.Context, id string) error {
	return a.ceHandler.ContainerStart(ctx, id)
}

func (a *Api) StopContainer(_ context.Context, id string) (string, error) {
	return a.jobHandler.Create(fmt.Sprintf("stop container '%s'", id), func(ctx context.Context, cf context.CancelFunc) error {
		defer cf()
		err := a.ceHandler.ContainerStop(ctx, id)
		if err == nil {
			err = ctx.Err()
		}
		return err
	})
}

func (a *Api) RestartContainer(_ context.Context, id string) (string, error) {
	return a.jobHandler.Create(fmt.Sprintf("restart container '%s'", id), func(ctx context.Context, cf context.CancelFunc) error {
		defer cf()
		err := a.ceHandler.ContainerRestart(ctx, id)
		if err == nil {
			err = ctx.Err()
		}
		return err
	})
}

func (a *Api) RemoveContainer(ctx context.Context, id string) error {
	return a.ceHandler.ContainerRemove(ctx, id)
}

func (a *Api) GetContainerLog(ctx context.Context, id string, logOptions model.LogOptions) (io.ReadCloser, error) {
	return a.ceHandler.ContainerLog(ctx, id, logOptions)
}
