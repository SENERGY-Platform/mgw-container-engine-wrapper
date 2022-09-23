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
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

func (d Docker) ListNetworks(ctx context.Context, filter [][2]string) ([]itf.Network, error) {
	var n []itf.Network
	nr, err := d.client.NetworkList(ctx, types.NetworkListOptions{Filters: util.GenFilterArgs(filter)})
	if err != nil {
		return n, err
	}
	for _, r := range nr {
		if nType, ok := util.NetTypeMap[r.Driver]; ok {
			s, gw := util.ParseNetIPAMConfig(r.IPAM.Config)
			n = append(n, itf.Network{
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

func (d Docker) NetworkInfo(ctx context.Context, id string) (itf.Network, error) {
	n := itf.Network{}
	nr, err := d.client.NetworkInspect(ctx, id, types.NetworkInspectOptions{})
	if err != nil {
		return n, err
	}
	s, gw := util.ParseNetIPAMConfig(nr.IPAM.Config)
	n.ID = nr.ID
	n.Name = nr.Name
	n.Type = util.NetTypeMap[nr.Driver]
	n.Subnet = s
	n.Gateway = gw
	return n, nil
}

func (d Docker) NetworkCreate(ctx context.Context, net itf.Network) error {
	res, err := d.client.NetworkCreate(ctx, net.Name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         util.NetTypeRMap[net.Type],
		Attachable:     true,
		IPAM: &network.IPAM{
			Config: util.GenNetIPAMConfig(net),
		},
	})
	if err != nil {
		return err
	}
	if res.Warning != "" {
		dmUtil.Logger.Warning(res.Warning)
	}
	return nil
}

func (d Docker) NetworkRemove(ctx context.Context, id string) error {
	return d.client.NetworkRemove(ctx, id)
}
