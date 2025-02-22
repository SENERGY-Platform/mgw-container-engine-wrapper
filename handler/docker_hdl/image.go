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
	hdl_util "github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker_hdl/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"strings"
)

func (h *Handler) ListImages(ctx context.Context, filter model.ImageFilter) ([]model.Image, error) {
	var images []model.Image
	il, err := h.client.ImageList(ctx, image.ListOptions{Filters: hdl_util.GenImageFilterArgs(filter)})
	if err != nil {
		return images, model.NewInternalError(err)
	}
	for _, is := range il {
		if filter.Name != "" && !inTags(is.RepoTags, filter.Name, filter.Tag) {
			continue
		}
		img := model.Image{
			ID: is.ID,
			//Created: is.Created,
			Size:    is.Size,
			Tags:    is.RepoTags,
			Digests: is.RepoDigests,
			Labels:  is.Labels,
		}
		if i, _, err := h.client.ImageInspectWithRaw(ctx, is.ID); err != nil {
			util.Logger.Errorf("inspecting image '%s' failed: %s", is.ID, err)
		} else {
			if ti, err := hdl_util.ParseTimestamp(i.Created); err != nil {
				util.Logger.Errorf("parsing created timestamp for image '%s' failed: %s", is.ID, err)
			} else {
				img.Created = ti.UTC()
			}
			img.Arch = i.Architecture
		}
		images = append(images, img)
	}
	return images, nil
}

func (h *Handler) ImageInfo(ctx context.Context, id string) (model.Image, error) {
	img := model.Image{}
	i, _, err := h.client.ImageInspectWithRaw(ctx, id)
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
	if ti, err := hdl_util.ParseTimestamp(i.Created); err != nil {
		util.Logger.Errorf("parsing created timestamp for image '%s' failed: %s", i.ID, err)
	} else {
		img.Created = ti.UTC()
	}
	return img, nil
}

func (h *Handler) ImagePull(ctx context.Context, id string) error {
	rc, err := h.client.ImagePull(ctx, id, image.PullOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	defer rc.Close()
	jd := json.NewDecoder(rc)
	var msg hdl_util.ImgPullResp
	for {
		if err := jd.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			} else {
				return model.NewInternalError(err)
			}
		}
		util.Logger.Debugf("pulling image '%s': %s", id, msg)
	}
	if msg.Message != "" {
		return model.NewInternalError(errors.New(msg.Message))
	}
	return nil
}

func (h *Handler) ImageRemove(ctx context.Context, id string) error {
	if _, err := h.client.ImageRemove(ctx, id, image.RemoveOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}

func (h *Handler) PruneImages(ctx context.Context) error {
	_, err := h.client.ImagesPrune(ctx, filters.Args{})
	return err
}

func inTags(list []string, name, tag string) bool {
	if tag != "" {
		name += ":" + tag
		for _, item := range list {
			if item == name {
				return true
			}
		}
	} else {
		for _, item := range list {
			if strings.HasPrefix(item, name) {
				return true
			}
		}
	}
	return false
}
