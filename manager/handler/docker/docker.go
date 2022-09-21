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
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
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
	info["api"] = d.client.ClientVersion()
	return info, nil
}

func (d *Docker) Close() error {
	return d.client.Close()
}

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

func (d *Docker) ListContainers(ctx context.Context, filter [][2]string) ([]itf.Container, error) {
	if cl, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: genFilterArgs(filter)}); err != nil {
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
				if tc, err := parseTimestamp(ci.Created); err != nil {
					util.Logger.Error(err)
				} else {
					ctr.Created = tc
				}
				if c.State == "running" {
					if ts, err := parseTimestamp(ci.State.StartedAt); err != nil {
						util.Logger.Error(err)
					} else {
						ctr.Started = &ts
					}
				}
				ctr.Image = ci.Config.Image
				ctr.EnvVars = parseEnv(ci.Config.Env)
				ctr.Ports = parsePortSetAndMap(ci.Config.ExposedPorts, ci.NetworkSettings.Ports)
				if len(ci.HostConfig.Mounts) > 0 {
					ctr.Mounts = parseMounts(ci.HostConfig.Mounts)
				}
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
		var mts []itf.Mount
		if len(c.HostConfig.Mounts) > 0 {
			mts = parseMounts(c.HostConfig.Mounts)
		} else {
			mts = parseMountPoints(c.Mounts)
		}
		ctr = itf.Container{
			ID:       c.ID,
			Name:     c.Name,
			State:    stateMap[c.State.Status],
			Image:    c.Config.Image,
			ImageID:  c.Image,
			EnvVars:  parseEnv(c.Config.Env),
			Labels:   c.Config.Labels,
			Mounts:   mts,
			Ports:    parsePortSetAndMap(c.Config.ExposedPorts, c.NetworkSettings.Ports),
			Networks: parseEndpointSettings(c.NetworkSettings.Networks),
			RunConfig: itf.RunConfig{
				RestartStrategy: restartPolicyMap[c.HostConfig.RestartPolicy.Name],
				Retries:         c.HostConfig.RestartPolicy.MaximumRetryCount,
				RemoveAfterRun:  c.HostConfig.AutoRemove,
				StopTimeout:     parseStopTimeout(c.Config.StopTimeout),
			},
		}
		if tc, err := parseTimestamp(c.Created); err != nil {
			util.Logger.Error(err)
		} else {
			ctr.Created = tc
		}
		if c.State.Status == "running" {
			if ts, err := parseTimestamp(c.State.StartedAt); err != nil {
				util.Logger.Error(err)
			} else {
				ctr.Started = &ts
			}
		}
	}
	return ctr, nil
}

func (d *Docker) ContainerCreate(ctx context.Context, ctrConf itf.Container) (string, error) {
	cConfig := &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Env:          genEnv(ctrConf.EnvVars),
		Image:        ctrConf.Image,
		Labels:       ctrConf.Labels,
		StopTimeout:  genStopTimeout(ctrConf.RunConfig.StopTimeout),
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
				err2 := d.ContainerRemove(ctx, res.ID)
				if err2 != nil {
					util.Logger.Error(err2)
				}
				return "", err
			}
		}
	}
	return res.ID, nil
}

func (d *Docker) ContainerRemove(ctx context.Context, id string) error {
	return d.client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
	})
}

func (d *Docker) ContainerStart(ctx context.Context, id string) error {
	return d.client.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

func (d *Docker) ContainerStop(ctx context.Context, id string) error {
	return d.client.ContainerStop(ctx, id, nil)
}

func (d *Docker) ContainerRestart(ctx context.Context, id string) error {
	return d.client.ContainerRestart(ctx, id, nil)
}

func (d *Docker) ContainerLog(ctx context.Context, id string) (io.ReadCloser, error) {
	var lr LogReader
	rc, err := d.client.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      "",
		Until:      "",
		Timestamps: false,
		Follow:     false,
		Tail:       "", // num of lines
		Details:    true,
	})
	if err != nil {
		return &lr, err
	}
	return NewLogReader(rc), nil
}

func (d *Docker) ListImages(ctx context.Context, filter [][2]string) ([]itf.Image, error) {
	if il, err := d.client.ImageList(ctx, types.ImageListOptions{All: true, Filters: genFilterArgs(filter)}); err != nil {
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
				util.Logger.Error(err)
			} else {
				if ti, err := parseTimestamp(i.Created); err != nil {
					util.Logger.Error(err)
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
		if ti, err := parseTimestamp(i.Created); err != nil {
			util.Logger.Error(err)
		} else {
			img.Created = ti
		}
	}
	return img, nil
}

func (d *Docker) ImagePull(ctx context.Context, id string) error {
	if res, err := d.client.ImagePull(ctx, id, types.ImagePullOptions{}); err != nil {
		return err
	} else {
		defer res.Close()
		jd := json.NewDecoder(res)
		var msg ImgPullResp
		for {
			if err := jd.Decode(&msg); err != nil {
				if err == io.EOF {
					break
				} else {
					return err
				}
			}
			util.Logger.Debug(msg)
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
		util.Logger.Debug(res)
	}
	return nil
}

func (d *Docker) PruneImages(ctx context.Context) error {
	if res, err := d.client.ImagesPrune(ctx, filters.Args{}); err != nil {
		return err
	} else {
		util.Logger.Debug(res)
	}
	return nil
}
