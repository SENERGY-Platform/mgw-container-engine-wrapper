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

package wrapper

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"io"
)

func (a *Wrapper) GetContainers(ctx context.Context, filter model.ContainerFilter) ([]model.Container, error) {
	return a.ceHandler.ListContainers(ctx, filter)
}

func (a *Wrapper) GetContainer(ctx context.Context, id string) (model.Container, error) {
	return a.ceHandler.ContainerInfo(ctx, id)
}

func (a *Wrapper) CreateContainer(ctx context.Context, container model.Container) (string, error) {
	return a.ceHandler.ContainerCreate(ctx, container)
}

func (a *Wrapper) StartContainer(ctx context.Context, id string) error {
	return a.ceHandler.ContainerStart(ctx, id)
}

func (a *Wrapper) StopContainer(ctx context.Context, id string) (string, error) {
	return a.jobHandler.Create(ctx, fmt.Sprintf("stop container '%s'", id), func(ctx context.Context, cf context.CancelFunc) (any, error) {
		defer cf()
		err := a.ceHandler.ContainerStop(ctx, id)
		if err == nil {
			err = ctx.Err()
		}
		return nil, err
	})
}

func (a *Wrapper) RestartContainer(ctx context.Context, id string) (string, error) {
	return a.jobHandler.Create(ctx, fmt.Sprintf("restart container '%s'", id), func(ctx context.Context, cf context.CancelFunc) (any, error) {
		defer cf()
		err := a.ceHandler.ContainerRestart(ctx, id)
		if err == nil {
			err = ctx.Err()
		}
		return nil, err
	})
}

func (a *Wrapper) RemoveContainer(ctx context.Context, id string, force bool) error {
	return a.ceHandler.ContainerRemove(ctx, id, force)
}

func (a *Wrapper) GetContainerLog(ctx context.Context, id string, logOptions model.LogFilter) (io.ReadCloser, error) {
	return a.ceHandler.ContainerLog(ctx, id, logOptions)
}

func (a *Wrapper) ContainerExec(ctx context.Context, id string, exeConf model.ExecConfig) (string, error) {
	return a.jobHandler.Create(ctx, fmt.Sprintf("container execute '%+v'", exeConf), func(ctx context.Context, cf context.CancelFunc) (any, error) {
		defer cf()
		err := a.ceHandler.ContainerExec(ctx, id, exeConf)
		if err == nil {
			err = ctx.Err()
		}
		return nil, err
	})
}
