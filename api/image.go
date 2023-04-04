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

package api

import (
	"container-engine-wrapper/model"
	"context"
	"fmt"
)

func (a *Api) GetImages(ctx context.Context, filter model.ImageFilter) ([]model.Image, error) {
	return a.ceHandler.ListImages(ctx, filter)
}

func (a *Api) GetImage(ctx context.Context, id string) (model.Image, error) {
	return a.ceHandler.ImageInfo(ctx, id)
}

func (a *Api) AddImage(_ context.Context, img string) (string, error) {
	return a.jobHandler.Create(fmt.Sprintf("add image '%s'", img), func(ctx context.Context, cf context.CancelFunc) error {
		defer cf()
		err := a.ceHandler.ImagePull(ctx, img)
		if err == nil {
			err = ctx.Err()
		}
		return err
	})
}

func (a *Api) RemoveImage(ctx context.Context, id string) error {
	return a.ceHandler.ImageRemove(ctx, id)
}
