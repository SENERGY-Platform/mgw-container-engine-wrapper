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

func (d *Docker) ContainerInfo(ctx context.Context, id string) (*itf.Container, error) {
	if c, err := d.client.ContainerInspect(ctx, id); err != nil {
		return nil, err
	} else {
		return &itf.Container{
			ID:          c.ID,
			Name:        c.Name,
			ImageID:     c.Image,
			Image:       c.Config.Image,
			Status:      c.State.Status,
			StartPolicy: c.HostConfig.RestartPolicy.Name,
			CreatedAt:   c.Created,
			StartedAt:   c.State.StartedAt,
			Hostname:    c.Config.Hostname,
			Env:         parseContainerEnvVars(c.Config.Env),
			Mounts:      parseContainerMounts(c.Mounts),
			Labels:      c.Config.Labels,
			Ports:       parseContainerPorts(c.Config.ExposedPorts, c.NetworkSettings.Ports),
			Networks:    parseContainerNetworks(c.NetworkSettings.Networks),
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
