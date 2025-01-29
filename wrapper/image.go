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
)

func (a *Wrapper) GetImages(ctx context.Context, filter model.ImageFilter) ([]model.Image, error) {
	return a.ceHandler.ListImages(ctx, filter)
}

func (a *Wrapper) GetImage(ctx context.Context, id string) (model.Image, error) {
	return a.ceHandler.ImageInfo(ctx, id)
}

func (a *Wrapper) AddImage(ctx context.Context, img string) (string, error) {
	return a.jobHandler.Create(ctx, fmt.Sprintf("add image '%s'", img), func(ctx context.Context, cf context.CancelFunc) (any, error) {
		defer cf()
		err := a.ceHandler.ImagePull(ctx, img)
		if err == nil {
			err = ctx.Err()
		}
		return nil, err
	})
}

func (a *Wrapper) RemoveImage(ctx context.Context, id string) error {
	return a.ceHandler.ImageRemove(ctx, id)
}
