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
)

func (a *Api) GetImages(ctx context.Context, filter itf.ImageFilter) ([]model.Image, error) {
	return a.ceHandler.ListImages(ctx, filter)
}

func (a *Api) GetImage(ctx context.Context, id string) (model.Image, error) {
	return a.ceHandler.ImageInfo(ctx, id)
}

func (a *Api) AddImage(_ context.Context, img string) (string, error) {
	jId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	ctx, cf := context.WithCancel(a.jobHandler.Context())
	j := itf.NewJob(ctx, cf, jId.String(), model.JobOrgRequest{
		Method: gc.Request.Method,
		Uri:    gc.Request.RequestURI,
		Body:   req,
	})
	j.SetTarget(func() {
		defer cf()
		e := a.ceHandler.ImagePull(ctx, img)
		if e == nil {
			e = ctx.Err()
		}
		j.SetError(e)
	})
	err = a.jobHandler.Add(jId.String(), j)
	if err != nil {
		return "", err
	}
	return jId.String(), nil
}

func (a *Api) RemoveImage(ctx context.Context, id string) error {
	return a.ceHandler.ImageRemove(ctx, id)
}
