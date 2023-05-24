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

package docker_hdl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker_hdl/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
)

func (d *Docker) ListImages(ctx context.Context, filter model.ImageFilter) ([]model.Image, error) {
	var images []model.Image
	il, err := d.client.ImageList(ctx, types.ImageListOptions{Filters: util.GenImageFilterArgs(filter)})
	if err != nil {
		return images, model.NewInternalError(err)
	}
	for _, is := range il {
		img := model.Image{
			ID: is.ID,
			//Created: is.Created,
			Size:    is.Size,
			Tags:    is.RepoTags,
			Digests: is.RepoDigests,
			Labels:  is.Labels,
		}
		if i, _, err := d.client.ImageInspectWithRaw(ctx, is.ID); err != nil {
			srv_base.Logger.Errorf("inspecting image '%s' failed: %s", is.ID, err)
		} else {
			if ti, err := util.ParseTimestamp(i.Created); err != nil {
				srv_base.Logger.Errorf("parsing created timestamp for image '%s' failed: %s", is.ID, err)
			} else {
				img.Created = ti.UTC()
			}
			img.Arch = i.Architecture
		}
		images = append(images, img)
	}
	return images, nil
}

func (d *Docker) ImageInfo(ctx context.Context, id string) (model.Image, error) {
	img := model.Image{}
	i, _, err := d.client.ImageInspectWithRaw(ctx, id)
	if err != nil {
		if client.IsErrNotFound(err) {
			return model.Image{}, model.NewNotFoundError(err)
		}
		return model.Image{}, model.NewInternalError(err)
	}
	img.ID = i.ID
	img.Size = i.Size
	img.Arch = i.Architecture
	img.Tags = i.RepoTags
	img.Digests = i.RepoDigests
	img.Labels = i.Config.Labels
	if ti, err := util.ParseTimestamp(i.Created); err != nil {
		srv_base.Logger.Errorf("parsing created timestamp for image '%s' failed: %s", i.ID, err)
	} else {
		img.Created = ti.UTC()
	}
	return img, nil
}

func (d *Docker) ImagePull(ctx context.Context, id string) error {
	rc, err := d.client.ImagePull(ctx, id, types.ImagePullOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	defer rc.Close()
	jd := json.NewDecoder(rc)
	var msg util.ImgPullResp
	for {
		if err := jd.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			} else {
				return model.NewInternalError(err)
			}
		}
		srv_base.Logger.Debugf("pulling image '%s': %s", id, msg)
	}
	if msg.Message != "" {
		return model.NewInternalError(errors.New(msg.Message))
	}
	return nil
}

func (d *Docker) ImageRemove(ctx context.Context, id string) error {
	if _, err := d.client.ImageRemove(ctx, id, types.ImageRemoveOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}

func (d *Docker) PruneImages(ctx context.Context) error {
	_, err := d.client.ImagesPrune(ctx, filters.Args{})
	return err
}
