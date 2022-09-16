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
	"deployment-manager/manager/ce-handler/itf"
	"deployment-manager/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Docker struct {
	client *client.Client
}

func New() (*Docker, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Docker{client: c}, nil
}

func (d *Docker) ServerInfo(ctx context.Context) (map[string]string, error) {
	srvVer, err := d.client.ServerVersion(ctx)
	if err != nil {
		return nil, err
	}
	info := make(map[string]string, len(srvVer.Components))
	for i := 0; i < len(srvVer.Components); i++ {
		info[srvVer.Components[i].Name] = srvVer.Components[i].Version
	}
	return info, nil
}

func (d *Docker) ListContainers(ctx context.Context, filter [][2]string) (map[string]*itf.Container, error) {
	var f filters.Args
	if filter != nil && len(filter) > 0 {
		f = filters.NewArgs()
		for _, i := range filter {
			f.Add(i[0], i[1])
		}
	}
	if cl, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: f}); err != nil {
		return nil, err
	} else {
		cm := make(map[string]*itf.Container, len(cl))
		for _, c := range cl {
			ctr := &itf.Container{
				ID:      c.ID,
				ImageID: c.ImageID,
				State:   stateMap[c.State],
				Mounts:  parseContainerMounts(c.Mounts),
				Labels:  c.Labels,
			}
			if ci, err := d.client.ContainerInspect(ctx, c.ID); err != nil {
				util.Logger.Error(err)
			} else {
				ctr.Name = ci.Name
				ctr.Image = ci.Config.Image
				ctr.RestartConfig = parseRestartPolicy(ci.HostConfig.RestartPolicy)
				ctr.Created = ci.Created
				ctr.Started = ci.State.StartedAt
				ctr.Env = parseContainerEnvVars(ci.Config.Env)
				ctr.Networks = parseContainerNetworks(ci.NetworkSettings.Networks)
				ctr.Ports = parseContainerPorts(ci.Config.ExposedPorts, ci.NetworkSettings.Ports)
			}
			cm[c.ID] = ctr
		}
		return cm, nil
	}
}

func (d *Docker) ContainerInfo(ctx context.Context, id string) (*itf.Container, error) {
	if c, err := d.client.ContainerInspect(ctx, id); err != nil {
		return nil, err
	} else {
		return &itf.Container{
			ID:            c.ID,
			Name:          c.Name,
			Image:         c.Config.Image,
			ImageID:       c.Image,
			State:         stateMap[c.State.Status],
			RestartConfig: parseRestartPolicy(c.HostConfig.RestartPolicy),
			Created:       c.Created,
			Started:       c.State.StartedAt,
			Env:           parseContainerEnvVars(c.Config.Env),
			Mounts:        parseContainerMounts(c.Mounts),
			Labels:        c.Config.Labels,
			Ports:         parseContainerPorts(c.Config.ExposedPorts, c.NetworkSettings.Ports),
			Networks:      parseContainerNetworks(c.NetworkSettings.Networks),
		}, nil
	}
}

func (d *Docker) ImageInfo(ctx context.Context, id string) (*itf.Image, error) {
	if i, _, err := d.client.ImageInspectWithRaw(ctx, id); err != nil {
		return nil, err
	} else {
		return &itf.Image{
			ID:      i.ID,
			Created: i.Created,
			Size:    i.Size,
			Arch:    i.Architecture,
			Tags:    i.RepoTags,
			Digests: i.RepoDigests,
		}, nil
	}
}
