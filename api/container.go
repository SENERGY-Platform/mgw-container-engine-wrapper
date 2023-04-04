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
	"container-engine-wrapper/itf"
	"container-engine-wrapper/model"
	"context"
	"github.com/google/uuid"
	"io"
)

func (a *Api) GetContainers(ctx context.Context, filter itf.ContainerFilter) ([]model.Container, error) {
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
	return a.ctrlCtrAsJob(id, a.ceHandler.ContainerStop)
}

func (a *Api) RestartContainer(_ context.Context, id string) (string, error) {
	return a.ctrlCtrAsJob(id, a.ceHandler.ContainerRestart)
}

func (a *Api) RemoveContainer(ctx context.Context, id string) error {
	return a.ceHandler.ContainerRemove(ctx, id)
}

func (a *Api) GetContainerLog(ctx context.Context, id string, logOptions itf.LogOptions) (io.ReadCloser, error) {
	return a.ceHandler.ContainerLog(ctx, id, logOptions)
}

func (a *Api) ctrlCtrAsJob(id string, f func(context.Context, string) error) (string, error) {
	jId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	ctx, cf := context.WithCancel(a.jobHandler.Context())
	j := itf.NewJob(ctx, cf, jId.String(), model.JobOrgRequest{
		Method: gc.Request.Method,
		Uri:    gc.Request.RequestURI,
	})
	j.SetTarget(func() {
		defer cf()
		e := f(ctx, id)
		if e == nil {
			e = ctx.Err()
		}
		j.SetError(e)
	})
	err = a.jobHandler.Add(jId.String(), j)
	if err != nil {
		return "", err
	}
	return jId.String(), err
}
