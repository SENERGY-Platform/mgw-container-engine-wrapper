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
	"deployment-manager/manager/itf"
	"deployment-manager/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

func (d *Docker) ListNetworks(ctx context.Context, filter [][2]string) ([]itf.Network, error) {
	if nr, err := d.client.NetworkList(ctx, types.NetworkListOptions{Filters: genFilterArgs(filter)}); err != nil {
		return nil, err
	} else {
		var n []itf.Network
		for _, r := range nr {
			if nType, ok := netTypeMap[r.Driver]; ok {
				s, gw := parseNetIPAMConfig(r.IPAM.Config)
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
}

func (d *Docker) NetworkInfo(ctx context.Context, id string) (itf.Network, error) {
	var n itf.Network
	if nr, err := d.client.NetworkInspect(ctx, id, types.NetworkInspectOptions{}); err != nil {
		return n, err
	} else {
		s, gw := parseNetIPAMConfig(nr.IPAM.Config)
		n = itf.Network{
			ID:      nr.ID,
			Name:    nr.Name,
			Type:    netTypeMap[nr.Driver],
			Subnet:  s,
			Gateway: gw,
		}
	}
	return n, nil
}

func (d *Docker) NetworkCreate(ctx context.Context, net itf.Network) error {
	if res, err := d.client.NetworkCreate(ctx, net.Name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         netTypeRMap[net.Type],
		Attachable:     true,
		IPAM: &network.IPAM{
			Config: genNetIPAMConfig(net),
		},
	}); err != nil {
		return err
	} else {
		util.Logger.Debug(res)
	}
	return nil
}

func (d *Docker) NetworkRemove(ctx context.Context, id string) error {
	return d.client.NetworkRemove(ctx, id)
}
