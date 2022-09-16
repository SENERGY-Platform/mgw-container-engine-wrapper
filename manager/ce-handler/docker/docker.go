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
				State:   stateMap[c.State],
				ImageID: c.ImageID,
				Labels:  c.Labels,
				Mounts:  parseContainerMounts(c.Mounts),
			}
			if ci, err := d.client.ContainerInspect(ctx, c.ID); err != nil {
				util.Logger.Error(err)
			} else {

				ctr.Name = ci.Name
				ctr.Created = ci.Created
				ctr.Started = ci.State.StartedAt
				ctr.Image = ci.Config.Image
				ctr.EnvVars = parseContainerEnvVars(ci.Config.Env)
				ctr.Ports = parseContainerPorts(ci.Config.ExposedPorts, ci.NetworkSettings.Ports)
				ctr.Networks = parseContainerNetworks(ci.NetworkSettings.Networks)
				ctr.RunConfig = itf.RunConfig{
					RestartStrategy: restartPolicyMap[ci.HostConfig.RestartPolicy.Name],
					Retries:         ci.HostConfig.RestartPolicy.MaximumRetryCount,
					RemoveAfterRun:  ci.HostConfig.AutoRemove,
					StopTimeout:     parseStopTimeout(ci.Config.StopTimeout),
				}
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
			ID:       c.ID,
			Name:     c.Name,
			State:    stateMap[c.State.Status],
			Created:  c.Created,
			Started:  c.State.StartedAt,
			Image:    c.Config.Image,
			ImageID:  c.Image,
			EnvVars:  parseContainerEnvVars(c.Config.Env),
			Labels:   c.Config.Labels,
			Mounts:   parseContainerMounts(c.Mounts),
			Ports:    parseContainerPorts(c.Config.ExposedPorts, c.NetworkSettings.Ports),
			Networks: parseContainerNetworks(c.NetworkSettings.Networks),
			RunConfig: itf.RunConfig{
				RestartStrategy: restartPolicyMap[c.HostConfig.RestartPolicy.Name],
				Retries:         c.HostConfig.RestartPolicy.MaximumRetryCount,
				RemoveAfterRun:  c.HostConfig.AutoRemove,
				StopTimeout:     parseStopTimeout(c.Config.StopTimeout),
			},
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
