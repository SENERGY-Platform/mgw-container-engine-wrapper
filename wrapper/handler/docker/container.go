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
	"container-engine-wrapper/wrapper/handler/docker/util"
	"context"
	"fmt"
	"github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/SENERGY-Platform/go-service-base/srv-base/types"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/wrapper/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (d *Docker) ListContainers(ctx context.Context, filter model.ContainerFilter) ([]model.Container, error) {
	var csl []model.Container
	cl, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: util.GenContainerFilterArgs(filter)})
	if err != nil {
		return csl, srv_base_types.NewError(http.StatusInternalServerError, "listing containers failed", err)
	}
	for _, c := range cl {
		ctr := model.Container{
			ID:       c.ID,
			State:    util.GetConst(c.State, util.StateMap),
			ImageID:  c.ImageID,
			Labels:   c.Labels,
			Mounts:   util.ParseMountPoints(c.Mounts),
			Networks: util.ParseEndpointSettings(c.NetworkSettings.Networks),
		}
		if ci, err := d.client.ContainerInspect(ctx, c.ID); err != nil {
			srv_base.Logger.Errorf("inspecting container '%s' failed: %s", c.ID, err)
		} else {
			ctr.Name = util.ParseContainerName(ci.Name)
			if tc, err := util.ParseTimestamp(ci.Created); err != nil {
				srv_base.Logger.Errorf("parsing created timestamp for container '%s' failed: %s", c.ID, err)
			} else {
				ctr.Created = tc
			}
			if c.State == "running" {
				if ts, err := util.ParseTimestamp(ci.State.StartedAt); err != nil {
					srv_base.Logger.Errorf("parsing started timestamp for container '%s' failed: %s", c.ID, err)
				} else {
					ctr.Started = &ts
				}
			}
			ctr.Image = ci.Config.Image
			ctr.EnvVars = util.ParseEnv(ci.Config.Env)
			if ports, err := util.ParsePortSetAndMap(ci.Config.ExposedPorts, ci.NetworkSettings.Ports); err != nil {
				srv_base.Logger.Errorf("parsing ports for container '%s' failed: %s", c.ID, err)
			} else {
				ctr.Ports = ports
			}
			if len(ci.HostConfig.Mounts) > 0 {
				ctr.Mounts = util.ParseMounts(ci.HostConfig.Mounts)
			}
			ctr.Networks = util.ParseEndpointSettings(ci.NetworkSettings.Networks)
			strategy, retries := util.ParseRestartPolicy(ci.HostConfig.RestartPolicy)
			ctr.RunConfig = model.RunConfig{
				RestartStrategy: strategy,
				Retries:         retries,
				RemoveAfterRun:  ci.HostConfig.AutoRemove,
				StopTimeout:     util.ParseStopTimeout(ci.Config.StopTimeout),
			}
			if ci.Config.StopSignal != "" {
				ctr.RunConfig.StopSignal = &ci.Config.StopSignal
			}
		}
		csl = append(csl, ctr)
	}
	return csl, nil
}

func (d *Docker) ContainerInfo(ctx context.Context, id string) (model.Container, error) {
	ctr := model.Container{}
	c, err := d.client.ContainerInspect(ctx, id)
	if err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return ctr, srv_base_types.NewError(code, fmt.Sprintf("retrieving info for container '%s' failed", id), err)
	}
	var mts []model.Mount
	if len(c.HostConfig.Mounts) > 0 {
		mts = util.ParseMounts(c.HostConfig.Mounts)
	} else {
		mts = util.ParseMountPoints(c.Mounts)
	}
	ctr.ID = c.ID
	ctr.Name = util.ParseContainerName(c.Name)
	ctr.State = util.GetConst(c.State.Status, util.StateMap)
	ctr.Image = c.Config.Image
	ctr.ImageID = c.Image
	ctr.EnvVars = util.ParseEnv(c.Config.Env)
	ctr.Labels = c.Config.Labels
	ctr.Mounts = mts
	if ports, err := util.ParsePortSetAndMap(c.Config.ExposedPorts, c.NetworkSettings.Ports); err != nil {
		srv_base.Logger.Errorf("parsing ports for container '%s' failed: %s", c.ID, err)
	} else {
		ctr.Ports = ports
	}
	ctr.Networks = util.ParseEndpointSettings(c.NetworkSettings.Networks)
	strategy, retries := util.ParseRestartPolicy(c.HostConfig.RestartPolicy)
	ctr.RunConfig = model.RunConfig{
		RestartStrategy: strategy,
		Retries:         retries,
		RemoveAfterRun:  c.HostConfig.AutoRemove,
		StopTimeout:     util.ParseStopTimeout(c.Config.StopTimeout),
	}
	if c.Config.StopSignal != "" {
		ctr.RunConfig.StopSignal = &c.Config.StopSignal
	}
	if tc, err := util.ParseTimestamp(c.Created); err != nil {
		srv_base.Logger.Errorf("parsing created timestamp for container '%s' failed: %s", c.ID, err)
	} else {
		ctr.Created = tc
	}
	if c.State.Status == "running" {
		if ts, err := util.ParseTimestamp(c.State.StartedAt); err != nil {
			srv_base.Logger.Errorf("parsing started timestamp for container '%s' failed: %s", c.ID, err)
		} else {
			ctr.Started = &ts
		}
	}
	return ctr, nil
}

func (d *Docker) ContainerCreate(ctx context.Context, ctrConf model.Container) (string, error) {
	cConfig := &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          ctrConf.RunConfig.PseudoTTY,
		Env:          util.GenEnv(ctrConf.EnvVars),
		Image:        ctrConf.Image,
		Labels:       ctrConf.Labels,
		StopTimeout:  util.GenStopTimeout(ctrConf.RunConfig.StopTimeout),
	}
	if ctrConf.RunConfig.StopSignal != nil {
		cConfig.StopSignal = *ctrConf.RunConfig.StopSignal
	}
	bindings, err := util.GenPortMap(ctrConf.Ports)
	if err != nil {
		return "", srv_base_types.NewError(http.StatusBadRequest, fmt.Sprintf("creating container '%s' failed", ctrConf.Name), err)
	}
	mts, err := util.GenMounts(ctrConf.Mounts)
	if err != nil {
		return "", srv_base_types.NewError(http.StatusBadRequest, fmt.Sprintf("creating container '%s' failed", ctrConf.Name), err)
	}
	rp, err := util.GenRestartPolicy(ctrConf.RunConfig.RestartStrategy, ctrConf.RunConfig.Retries)
	if err != nil {
		return "", srv_base_types.NewError(http.StatusBadRequest, fmt.Sprintf("creating container '%s' failed", ctrConf.Name), err)
	}
	hConfig := &container.HostConfig{
		PortBindings:  bindings,
		RestartPolicy: rp,
		AutoRemove:    ctrConf.RunConfig.RemoveAfterRun,
		Mounts:        mts,
	}
	err = util.CheckNetworks(ctrConf.Networks)
	if err != nil {
		return "", srv_base_types.NewError(http.StatusBadRequest, fmt.Sprintf("creating container '%s' failed", ctrConf.Name), err)
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
		return "", srv_base_types.NewError(http.StatusInternalServerError, fmt.Sprintf("creating container '%s' failed", ctrConf.Name), err)
	}
	if len(ctrConf.Networks) > 1 {
		for i := 1; i < len(ctrConf.Networks); i++ {
			err := d.client.NetworkConnect(ctx, ctrConf.Networks[i].Name, res.ID, &network.EndpointSettings{
				Aliases: ctrConf.Networks[i].DomainNames,
			})
			if err != nil {
				err2 := d.ContainerRemove(ctx, res.ID)
				if err2 != nil {
					srv_base.Logger.Errorf("removing container '%s' failed: %s", ctrConf.Name, err2)
				}
				return "", srv_base_types.NewError(http.StatusInternalServerError, fmt.Sprintf("creating container '%s' failed", ctrConf.Name), err)
			}
		}
	}
	if res.Warnings != nil && len(res.Warnings) > 0 {
		srv_base.Logger.Warningf("encountered warnings during creation of container '%s': %s", ctrConf.Name, res.Warnings)
	}
	return res.ID, nil
}

func (d *Docker) ContainerRemove(ctx context.Context, id string) error {
	if err := d.client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return srv_base_types.NewError(code, fmt.Sprintf("removing container '%s' failed", id), err)
	}
	return nil
}

func (d *Docker) ContainerStart(ctx context.Context, id string) error {
	if err := d.client.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return srv_base_types.NewError(code, fmt.Sprintf("starting container '%s' failed", id), err)
	}
	return nil
}

func (d *Docker) ContainerStop(ctx context.Context, id string) error {
	if err := d.client.ContainerStop(ctx, id, nil); err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return srv_base_types.NewError(code, fmt.Sprintf("stopping container '%s' failed", id), err)
	}
	return nil
}

func (d *Docker) ContainerRestart(ctx context.Context, id string) error {
	if err := d.client.ContainerRestart(ctx, id, nil); err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return srv_base_types.NewError(code, fmt.Sprintf("restarting container '%s' failed", id), err)
	}
	return nil
}

func (d *Docker) ContainerLog(ctx context.Context, id string, logOpt model.LogOptions) (io.ReadCloser, error) {
	clo := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	if logOpt.Since != nil {
		clo.Since = time.Time(*logOpt.Since).Format(time.RFC3339Nano)
	}
	if logOpt.Until != nil {
		clo.Until = time.Time(*logOpt.Until).Format(time.RFC3339Nano)
	}
	if logOpt.MaxLines > 0 {
		clo.Tail = strconv.FormatInt(int64(logOpt.MaxLines), 10)
	}
	rc, err := d.client.ContainerLogs(ctx, id, clo)
	if err != nil {
		code := http.StatusInternalServerError
		if client.IsErrNotFound(err) {
			code = http.StatusNotFound
		}
		return nil, srv_base_types.NewError(code, fmt.Sprintf("retrieving log for container '%s' failed", id), err)
	}
	return util.NewLogReader(rc), nil
}
