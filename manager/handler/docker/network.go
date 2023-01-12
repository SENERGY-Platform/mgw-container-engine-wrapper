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
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"net/http"
)

func (d *Docker) ListNetworks(ctx context.Context) ([]model.Network, error) {
	var n []model.Network
	nr, err := d.client.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return n, srv_base_types.NewError(http.StatusInternalServerError, "listing networks failed", err)
	}
	for _, r := range nr {
		if nType, ok := util.NetTypeMap[r.Driver]; ok {
			s, gw := util.ParseNetIPAMConfig(r.IPAM.Config)
			n = append(n, model.Network{
				ID:      r.ID,
				Name:    r.Name,
				Type:    nType,
				Subnet:  s,
				Gateway: gw,
			})
		}
	}
	return n, nil
}

func (d *Docker) NetworkInfo(ctx context.Context, id string) (model.Network, error) {
	n := model.Network{}
	nr, err := d.client.NetworkInspect(ctx, id, types.NetworkInspectOptions{})
	if err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return n, srv_base_types.NewError(code, fmt.Sprintf("retrieving info for network '%s' failed", id), err)
	}
	s, gw := util.ParseNetIPAMConfig(nr.IPAM.Config)
	n.ID = nr.ID
	n.Name = nr.Name
	n.Type = util.GetConst(nr.Driver, util.NetTypeMap)
	n.Subnet = s
	n.Gateway = gw
	return n, nil
}

func (d *Docker) NetworkCreate(ctx context.Context, net model.Network) error {
	if _, ok := model.NetworkTypeMap[net.Type]; !ok {
		return srv_base_types.NewError(http.StatusBadRequest, "", fmt.Errorf("invalid network type '%s'", net.Type))
	}
	res, err := d.client.NetworkCreate(ctx, net.Name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         util.NetTypeRMap[net.Type],
		Attachable:     true,
		IPAM: &network.IPAM{
			Config: util.GenNetIPAMConfig(net),
		},
	})
	if err != nil {
		return srv_base_types.NewError(http.StatusInternalServerError, fmt.Sprintf("creating network '%s' failed", net.Name), err)
	}
	if res.Warning != "" {
		srv_base.Logger.Warningf("encountered warnings during creation of network '%s': %s", net.Name, res.Warning)
	}
	return nil
}

func (d *Docker) NetworkRemove(ctx context.Context, id string) error {
	if err := d.client.NetworkRemove(ctx, id); err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return srv_base_types.NewError(code, fmt.Sprintf("removing network '%s' failed", id), err)
	}
	return nil
}
