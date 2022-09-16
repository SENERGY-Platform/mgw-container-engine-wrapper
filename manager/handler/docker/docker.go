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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
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

func (d *Docker) ListContainers(ctx context.Context, filter [][2]string) ([]itf.Container, error) {
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
		var csl []itf.Container
		for _, c := range cl {
			ctr := itf.Container{
				ID:       c.ID,
				State:    stateMap[c.State],
				ImageID:  c.ImageID,
				Labels:   c.Labels,
				Mounts:   parseMountPoints(c.Mounts),
				Networks: parseEndpointSettings(c.NetworkSettings.Networks),
			}
			if ci, err := d.client.ContainerInspect(ctx, c.ID); err != nil {
				util.Logger.Error(err)
			} else {
				ctr.Name = ci.Name
				ctr.Created = ci.Created
				ctr.Started = ci.State.StartedAt
				ctr.Image = ci.Config.Image
				ctr.EnvVars = parseEnv(ci.Config.Env)
				ctr.Ports = parsePortSetAndMap(ci.Config.ExposedPorts, ci.NetworkSettings.Ports)
				ctr.Networks = parseEndpointSettings(ci.NetworkSettings.Networks)
				ctr.RunConfig = itf.RunConfig{
					RestartStrategy: restartPolicyMap[ci.HostConfig.RestartPolicy.Name],
					Retries:         ci.HostConfig.RestartPolicy.MaximumRetryCount,
					RemoveAfterRun:  ci.HostConfig.AutoRemove,
					StopTimeout:     parseStopTimeout(ci.Config.StopTimeout),
				}
			}
			csl = append(csl, ctr)
		}
		return csl, nil
	}
}

func (d *Docker) ContainerInfo(ctx context.Context, id string) (itf.Container, error) {
	var ctr itf.Container
	if c, err := d.client.ContainerInspect(ctx, id); err != nil {
		return ctr, err
	} else {
		ctr = itf.Container{
			ID:       c.ID,
			Name:     c.Name,
			State:    stateMap[c.State.Status],
			Created:  c.Created,
			Started:  c.State.StartedAt,
			Image:    c.Config.Image,
			ImageID:  c.Image,
			EnvVars:  parseEnv(c.Config.Env),
			Labels:   c.Config.Labels,
			Mounts:   parseMountPoints(c.Mounts),
			Ports:    parsePortSetAndMap(c.Config.ExposedPorts, c.NetworkSettings.Ports),
			Networks: parseEndpointSettings(c.NetworkSettings.Networks),
			RunConfig: itf.RunConfig{
				RestartStrategy: restartPolicyMap[c.HostConfig.RestartPolicy.Name],
				Retries:         c.HostConfig.RestartPolicy.MaximumRetryCount,
				RemoveAfterRun:  c.HostConfig.AutoRemove,
				StopTimeout:     parseStopTimeout(c.Config.StopTimeout),
			},
		}
	}
	return ctr, nil
}

func (d *Docker) ContainerCreate(ctx context.Context, ctrConf itf.Container) (string, error) {
	cConfig := &container.Config{
		Env:         genEnv(ctrConf.EnvVars),
		Image:       ctrConf.Image,
		Labels:      ctrConf.Labels,
		StopTimeout: genStopTimeout(ctrConf.RunConfig.StopTimeout),
	}
	bindings, err := genPortMap(ctrConf.Ports)
	if err != nil {
		return "", err
	}
	mts, err := genMounts(ctrConf.Mounts)
	if err != nil {
		return "", err
	}
	hConfig := &container.HostConfig{
		PortBindings: bindings,
		RestartPolicy: container.RestartPolicy{
			Name:              restartPolicyRMap[ctrConf.RunConfig.RestartStrategy],
			MaximumRetryCount: ctrConf.RunConfig.Retries,
		},
		AutoRemove: ctrConf.RunConfig.RemoveAfterRun,
		Mounts:     mts,
	}
	err = checkNetworks(ctrConf.Networks)
	if err != nil {
		return "", err
	}
	var nConfig *network.NetworkingConfig
	if len(ctrConf.Networks) > 0 {
		nConfig = &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{
			ctrConf.Networks[0].Name: {
				Aliases: ctrConf.Networks[0].DomainNames,
			},
		}}
	}
	res, err := d.client.ContainerCreate(ctx, cConfig, hConfig, nConfig, nil, ctrConf.Name)
	if err != nil {
		return "", err
	}
	if res.Warnings != nil && len(res.Warnings) > 0 {
		util.Logger.Warning(res.Warnings)
	}
	if len(ctrConf.Networks) > 1 {
		for i := 1; i < len(ctrConf.Networks); i++ {
			err := d.client.NetworkConnect(ctx, ctrConf.Networks[i].Name, res.ID, &network.EndpointSettings{
				Aliases: ctrConf.Networks[i].DomainNames,
			})
			if err != nil {
				err2 := d.client.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{
					Force: true,
				})
				if err2 != nil {
					util.Logger.Error(err2)
				}
				return "", err
			}
		}
	}
	return res.ID, nil
}

func (d *Docker) ImageInfo(ctx context.Context, id string) (itf.Image, error) {
	var img itf.Image
	if i, _, err := d.client.ImageInspectWithRaw(ctx, id); err != nil {
		return img, err
	} else {
		img = itf.Image{
			ID:      i.ID,
			Created: i.Created,
			Size:    i.Size,
			Arch:    i.Architecture,
			Tags:    i.RepoTags,
			Digests: i.RepoDigests,
		}
	}
	return img, nil
}
