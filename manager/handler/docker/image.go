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

package docker

import (
	"context"
	"deployment-manager/manager/handler/docker/util"
	"deployment-manager/manager/itf"
	dmUtil "deployment-manager/util"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
)

func (d *Docker) ListImages(ctx context.Context, filter [][2]string) ([]itf.Image, error) {
	if il, err := d.client.ImageList(ctx, types.ImageListOptions{All: true, Filters: util.GenFilterArgs(filter)}); err != nil {
		return nil, err
	} else {
		var images []itf.Image
		for _, is := range il {
			img := itf.Image{
				ID: is.ID,
				//Created: is.Created,
				Size:    is.Size,
				Tags:    is.RepoTags,
				Digests: is.RepoDigests,
			}
			if i, _, err := d.client.ImageInspectWithRaw(ctx, is.ID); err != nil {
				dmUtil.Logger.Error(err)
			} else {
				if ti, err := util.ParseTimestamp(i.Created); err != nil {
					dmUtil.Logger.Error(err)
				} else {
					img.Created = ti
				}
				img.Arch = i.Architecture
			}
			images = append(images, img)
		}
		return images, nil
	}
}

func (d *Docker) ImageInfo(ctx context.Context, id string) (itf.Image, error) {
	var img itf.Image
	if i, _, err := d.client.ImageInspectWithRaw(ctx, id); err != nil {
		return img, err
	} else {
		img = itf.Image{
			ID:      i.ID,
			Size:    i.Size,
			Arch:    i.Architecture,
			Tags:    i.RepoTags,
			Digests: i.RepoDigests,
		}
		if ti, err := util.ParseTimestamp(i.Created); err != nil {
			dmUtil.Logger.Error(err)
		} else {
			img.Created = ti
		}
	}
	return img, nil
}

func (d *Docker) ImagePull(ctx context.Context, id string) error {
	if rc, err := d.client.ImagePull(ctx, id, types.ImagePullOptions{}); err != nil {
		return err
	} else {
		defer rc.Close()
		jd := json.NewDecoder(rc)
		var msg util.ImgPullResp
		for {
			if err := jd.Decode(&msg); err != nil {
				if err == io.EOF {
					break
				} else {
					return err
				}
			}
			dmUtil.Logger.Debug(msg)
		}
		if msg.Message != "" {
			return errors.New(msg.Message)
		}
	}
	return nil
}

func (d *Docker) ImageRemove(ctx context.Context, id string) error {
	if res, err := d.client.ImageRemove(ctx, id, types.ImageRemoveOptions{}); err != nil {
		return err
	} else {
		dmUtil.Logger.Debug(res)
	}
	return nil
}

func (d *Docker) PruneImages(ctx context.Context) error {
	if res, err := d.client.ImagesPrune(ctx, filters.Args{}); err != nil {
		return err
	} else {
		dmUtil.Logger.Debug(res)
	}
	return nil
}
