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
	"errors"
	hdl_util "github.com/SENERGY-Platform/mgw-container-engine-wrapper/handler/docker_hdl/util"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"strconv"
	"time"
)

func (d *Docker) ListContainers(ctx context.Context, filter model.ContainerFilter) ([]model.Container, error) {
	var csl []model.Container
	cl, err := d.client.ContainerList(ctx, container.ListOptions{All: true, Filters: hdl_util.GenContainerFilterArgs(filter)})
	if err != nil {
		return nil, model.NewInternalError(err)
	}
	for _, c := range cl {
		ctr := model.Container{
			ID:       c.ID,
			State:    hdl_util.GetConst(c.State, hdl_util.StateMap),
			ImageID:  c.ImageID,
			Labels:   c.Labels,
			Mounts:   hdl_util.ParseMountPoints(c.Mounts),
			Networks: hdl_util.ParseEndpointSettings(c.NetworkSettings.Networks),
		}
		if ci, err := d.client.ContainerInspect(ctx, c.ID); err != nil {
			util.Logger.Errorf("inspecting container '%s' failed: %s", c.ID, err)
		} else {
			ctr.Name = hdl_util.ParseContainerName(ci.Name)
			if tc, err := hdl_util.ParseTimestamp(ci.Created); err != nil {
				util.Logger.Errorf("parsing created timestamp for container '%s' failed: %s", c.ID, err)
			} else {
				ctr.Created = tc
			}
			if c.State == "running" {
				if ts, err := hdl_util.ParseTimestamp(ci.State.StartedAt); err != nil {
					util.Logger.Errorf("parsing started timestamp for container '%s' failed: %s", c.ID, err)
				} else {
					ctr.Started = &ts
				}
			}
			ctr.Image = ci.Config.Image
			ctr.EnvVars = hdl_util.ParseEnv(ci.Config.Env)
			if ports, err := hdl_util.ParsePortSetAndMap(ci.Config.ExposedPorts, ci.NetworkSettings.Ports); err != nil {
				util.Logger.Errorf("parsing ports for container '%s' failed: %s", c.ID, err)
			} else {
				ctr.Ports = ports
			}
			if len(ci.HostConfig.Mounts) > 0 {
				ctr.Mounts = hdl_util.ParseMounts(ci.HostConfig.Mounts)
			}
			ctr.Networks = hdl_util.ParseEndpointSettings(ci.NetworkSettings.Networks)
			strategy, retries := hdl_util.ParseRestartPolicy(ci.HostConfig.RestartPolicy)
			ctr.RunConfig = model.RunConfig{
				RestartStrategy: strategy,
				Retries:         retries,
				RemoveAfterRun:  ci.HostConfig.AutoRemove,
				StopTimeout:     hdl_util.ParseStopTimeout(ci.Config.StopTimeout),
			}
			if ci.Config.StopSignal != "" {
				ctr.RunConfig.StopSignal = &ci.Config.StopSignal
			}
			if ci.State.Health != nil {
				hs := hdl_util.GetConst(ci.State.Health.Status, hdl_util.HealthMap)
				ctr.Health = &hs
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
		if client.IsErrNotFound(err) {
			return model.Container{}, model.NewNotFoundError(err)
		}
		return model.Container{}, model.NewInternalError(err)
	}
	var mts []model.Mount
	if len(c.HostConfig.Mounts) > 0 {
		mts = hdl_util.ParseMounts(c.HostConfig.Mounts)
	} else {
		mts = hdl_util.ParseMountPoints(c.Mounts)
	}
	ctr.ID = c.ID
	ctr.Name = hdl_util.ParseContainerName(c.Name)
	ctr.State = hdl_util.GetConst(c.State.Status, hdl_util.StateMap)
	ctr.Image = c.Config.Image
	ctr.ImageID = c.Image
	ctr.EnvVars = hdl_util.ParseEnv(c.Config.Env)
	ctr.Labels = c.Config.Labels
	ctr.Mounts = mts
	if ports, err := hdl_util.ParsePortSetAndMap(c.Config.ExposedPorts, c.NetworkSettings.Ports); err != nil {
		util.Logger.Errorf("parsing ports for container '%s' failed: %s", c.ID, err)
	} else {
		ctr.Ports = ports
	}
	ctr.Networks = hdl_util.ParseEndpointSettings(c.NetworkSettings.Networks)
	strategy, retries := hdl_util.ParseRestartPolicy(c.HostConfig.RestartPolicy)
	ctr.RunConfig = model.RunConfig{
		RestartStrategy: strategy,
		Retries:         retries,
		RemoveAfterRun:  c.HostConfig.AutoRemove,
		StopTimeout:     hdl_util.ParseStopTimeout(c.Config.StopTimeout),
	}
	if c.Config.StopSignal != "" {
		ctr.RunConfig.StopSignal = &c.Config.StopSignal
	}
	if tc, err := hdl_util.ParseTimestamp(c.Created); err != nil {
		util.Logger.Errorf("parsing created timestamp for container '%s' failed: %s", c.ID, err)
	} else {
		ctr.Created = tc
	}
	if c.State.Status == "running" {
		if ts, err := hdl_util.ParseTimestamp(c.State.StartedAt); err != nil {
			util.Logger.Errorf("parsing started timestamp for container '%s' failed: %s", c.ID, err)
		} else {
			ctr.Started = &ts
		}
	}
	if c.State.Health != nil {
		hs := hdl_util.GetConst(c.State.Health.Status, hdl_util.HealthMap)
		ctr.Health = &hs
	}
	return ctr, nil
}

func (d *Docker) ContainerCreate(ctx context.Context, ctrConf model.Container) (string, error) {
	portMap, portSet, err := hdl_util.GenPortMap(ctrConf.Ports)
	if err != nil {
		return "", model.NewInvalidInputError(err)
	}
	cConfig := &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: portSet,
		Tty:          ctrConf.RunConfig.PseudoTTY,
		Env:          hdl_util.GenEnv(ctrConf.EnvVars),
		Image:        ctrConf.Image,
		Labels:       ctrConf.Labels,
		StopTimeout:  hdl_util.GenStopTimeout(ctrConf.RunConfig.StopTimeout),
	}
	if ctrConf.RunConfig.StopSignal != nil {
		cConfig.StopSignal = *ctrConf.RunConfig.StopSignal
	}
	if len(ctrConf.RunConfig.Command) > 0 {
		cConfig.Cmd = ctrConf.RunConfig.Command
	}
	mts, err := hdl_util.GenMounts(ctrConf.Mounts)
	if err != nil {
		return "", model.NewInvalidInputError(err)
	}
	dvs, err := hdl_util.GenDevices(ctrConf.Devices)
	rp, err := hdl_util.GenRestartPolicy(ctrConf.RunConfig.RestartStrategy, ctrConf.RunConfig.Retries)
	if err != nil {
		return "", model.NewInvalidInputError(err)
	}
	hConfig := &container.HostConfig{
		PortBindings:  portMap,
		RestartPolicy: rp,
		AutoRemove:    ctrConf.RunConfig.RemoveAfterRun,
		Mounts:        mts,
		Resources: container.Resources{
			Devices: dvs,
		},
	}
	err = hdl_util.CheckNetworks(ctrConf.Networks)
	if err != nil {
		return "", model.NewInvalidInputError(err)
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
		return "", model.NewInternalError(err)
	}
	if len(ctrConf.Networks) > 1 {
		for i := 1; i < len(ctrConf.Networks); i++ {
			err := d.client.NetworkConnect(ctx, ctrConf.Networks[i].Name, res.ID, &network.EndpointSettings{
				Aliases: ctrConf.Networks[i].DomainNames,
			})
			if err != nil {
				err2 := d.ContainerRemove(ctx, res.ID, true)
				if err2 != nil {
					util.Logger.Errorf("removing container '%s' failed: %s", ctrConf.Name, err2)
				}
				return "", model.NewInternalError(err)
			}
		}
	}
	if res.Warnings != nil && len(res.Warnings) > 0 {
		util.Logger.Warningf("encountered warnings during creation of container '%s': %s", ctrConf.Name, res.Warnings)
	}
	return res.ID, nil
}

func (d *Docker) ContainerRemove(ctx context.Context, id string, force bool) error {
	if err := d.client.ContainerRemove(ctx, id, container.RemoveOptions{Force: force}); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}

func (d *Docker) ContainerStart(ctx context.Context, id string) error {
	if err := d.client.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}

func (d *Docker) ContainerStop(ctx context.Context, id string) error {
	if err := d.client.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}

func (d *Docker) ContainerRestart(ctx context.Context, id string) error {
	if err := d.client.ContainerRestart(ctx, id, container.StopOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			return model.NewNotFoundError(err)
		}
		return model.NewInternalError(err)
	}
	return nil
}

func (d *Docker) ContainerLog(ctx context.Context, id string, logOpt model.LogFilter) (io.ReadCloser, error) {
	clo := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	if !logOpt.Since.IsZero() {
		clo.Since = logOpt.Since.Format(time.RFC3339Nano)
	}
	if !logOpt.Until.IsZero() {
		clo.Until = logOpt.Until.Format(time.RFC3339Nano)
	}
	if logOpt.MaxLines > 0 {
		clo.Tail = strconv.FormatInt(int64(logOpt.MaxLines), 10)
	}
	rc, err := d.client.ContainerLogs(ctx, id, clo)
	if err != nil {
		if client.IsErrNotFound(err) {
			return nil, model.NewNotFoundError(err)
		}
		return nil, model.NewInternalError(err)
	}
	return &hdl_util.RCWrapper{ReadCloser: rc}, nil
}

func (d *Docker) ContainerExec(ctx context.Context, id string, execOpt model.ExecConfig) error {
	eConf, err := d.client.ContainerExecCreate(ctx, id, types.ExecConfig{
		Tty:          execOpt.Tty,
		AttachStderr: true,
		AttachStdout: true,
		Env:          hdl_util.GenEnv(execOpt.EnvVars),
		WorkingDir:   execOpt.WorkDir,
		Cmd:          execOpt.Cmd,
	})
	if err != nil {
		return model.NewInternalError(err)
	}
	eAttach, err := d.client.ContainerExecAttach(ctx, eConf.ID, types.ExecStartCheck{Tty: execOpt.Tty})
	if err != nil {
		return model.NewInternalError(err)
	}
	defer eAttach.Close()
	eRes, err := d.awaitContainerExec(ctx, eConf.ID, time.Millisecond*250)
	if err != nil {
		return model.NewInternalError(err)
	}
	if eRes.ExitCode > 0 {
		bytes, err := io.ReadAll(eAttach.Reader)
		if err != nil {
			return model.NewInternalError(err)
		}
		return model.NewInternalError(errors.New(string(bytes)))
	}
	return nil
}

func (d *Docker) awaitContainerExec(ctx context.Context, execID string, delay time.Duration) (types.ContainerExecInspect, error) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return types.ContainerExecInspect{}, ctx.Err()
		case <-ticker.C:
			eRes, err := d.client.ContainerExecInspect(ctx, execID)
			if err != nil {
				return types.ContainerExecInspect{}, err
			}
			if !eRes.Running {
				return eRes, nil
			}
		}
	}
}
