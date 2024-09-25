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
	"fmt"
	hdl_util "github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker_hdl/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func (d *Docker) ListNetworks(ctx context.Context) ([]model.Network, error) {
	var n []model.Network
	nr, err := d.client.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, model.NewInternalError(err)
	}
	for _, r := range nr {
		if nType, ok := hdl_util.NetTypeMap[r.Driver]; ok {
			s, gw := hdl_util.ParseNetIPAMConfig(r.IPAM.Config)
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
	nr, err := d.client.NetworkInspect(ctx, id, network.InspectOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			return model.Network{}, model.NewNotFoundError(err)
		}
		return model.Network{}, model.NewInternalError(err)
	}
	s, gw := hdl_util.ParseNetIPAMConfig(nr.IPAM.Config)
	n.ID = nr.ID
	n.Name = nr.Name
	n.Type = hdl_util.GetConst(nr.Driver, hdl_util.NetTypeMap)
	n.Subnet = s
	n.Gateway = gw
	return n, nil
}

func (d *Docker) NetworkCreate(ctx context.Context, net model.Network) (string, error) {
	if _, ok := model.NetworkTypeMap[net.Type]; !ok {
		return "", model.NewInvalidInputError(fmt.Errorf("invalid network type '%s'", net.Type))
	}
	res, err := d.client.NetworkCreate(ctx, net.Name, network.CreateOptions{
		Driver:     hdl_util.NetTypeRMap[net.Type],
		Attachable: true,
		IPAM: &network.IPAM{
			Config: hdl_util.GenNetIPAMConfig(net),
		},
	})
	if err != nil {
		return "", model.NewInternalError(err)
	}
	if res.Warning != "" {
		util.Logger.Warningf("encountered warnings during creation of network '%s': %s", net.Name, res.Warning)
	}
	return res.ID, nil
}

func (d *Docker) NetworkRemove(ctx context.Context, id string) error {
	if err := d.client.NetworkRemove(ctx, id); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}
