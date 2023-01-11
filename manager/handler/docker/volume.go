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
	"container-engine-manager/manager/handler/docker/util"
	"context"
	"fmt"
	"github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/SENERGY-Platform/go-service-base/srv-base/types"
	"github.com/SENERGY-Platform/mgw-container-engine-manager/manager/model"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"net/http"
)

func (d *Docker) ListVolumes(ctx context.Context, filter model.VolumeFilter) ([]model.Volume, error) {
	var vols []model.Volume
	vls, err := d.client.VolumeList(ctx, util.GenVolumeFilterArgs(filter))
	if err != nil {
		return vols, srv_base_types.NewError(http.StatusInternalServerError, "listing volumes failed", err)
	}
	for _, vl := range vls.Volumes {
		vol := model.Volume{
			Name:   vl.Name,
			Labels: vl.Labels,
		}
		if ti, err := util.ParseTimestamp(vl.CreatedAt); err != nil {
			srv_base.Logger.Errorf("parsing created timestamp for volume '%s' failed: %s", vl.Name, err)
		} else {
			vol.Created = ti
		}
		vols = append(vols, vol)
	}
	return vols, nil
}

func (d *Docker) VolumeInfo(ctx context.Context, id string) (model.Volume, error) {
	vol := model.Volume{}
	vl, err := d.client.VolumeInspect(ctx, id)
	if err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return vol, srv_base_types.NewError(code, fmt.Sprintf("retrieving info for volume '%s' failed", id), err)
	}
	vol.Name = vl.Name
	vol.Labels = vl.Labels
	if ti, err := util.ParseTimestamp(vl.CreatedAt); err != nil {
		srv_base.Logger.Errorf("parsing created timestamp for volume '%s' failed: %s", vl.Name, err)
	} else {
		vol.Created = ti
	}
	return vol, nil
}

func (d *Docker) VolumeCreate(ctx context.Context, vol model.Volume) error {
	_, err := d.client.VolumeCreate(ctx, volume.VolumeCreateBody{Name: vol.Name, Labels: vol.Labels})
	if err != nil {
		return srv_base_types.NewError(http.StatusInternalServerError, fmt.Sprintf("creating volume '%s' failed", vol.Name), err)
	}
	return nil
}

func (d *Docker) VolumeRemove(ctx context.Context, id string) error {
	if err := d.client.VolumeRemove(ctx, id, false); err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return srv_base_types.NewError(code, fmt.Sprintf("removing volume '%s' failed", id), err)
	}
	return nil
}
