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
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/util"
	"github.com/docker/docker/client"
	"time"
)

type ContainerLogConf struct {
	Driver  string
	MaxSize string
	MaxFile int
}

type Docker struct {
	client     *client.Client
	ctrLogConf ContainerLogConf
}

func New(c *client.Client, ctrLogConf ContainerLogConf) (*Docker, error) {
	if ctrLogConf.Driver != "" && !isValidLoggingDriver(ctrLogConf.Driver) {
		return nil, errors.New("invalid logging driver: " + ctrLogConf.Driver)
	}
	return &Docker{
		client:     c,
		ctrLogConf: ctrLogConf,
	}, nil
}

func (d *Docker) ServerInfo(ctx context.Context, delay time.Duration) (map[string]string, error) {
	err := d.waitForServer(ctx, delay)
	if err != nil {
		return nil, err
	}
	info := map[string]string{}
	srvVer, err := d.client.ServerVersion(ctx)
	if err != nil {
		return info, err
	}
	for i := 0; i < len(srvVer.Components); i++ {
		info[srvVer.Components[i].Name] = srvVer.Components[i].Version
	}
	info["api"] = d.client.ClientVersion()
	return info, nil
}

func (d *Docker) waitForServer(ctx context.Context, delay time.Duration) error {
	_, err := d.client.Ping(ctx)
	if err == nil {
		return nil
	} else {
		if !client.IsErrConnectionFailed(err) {
			return err
		} else {
			util.Logger.Error(err)
		}
	}
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, err = d.client.Ping(ctx)
			if err == nil {
				return nil
			} else {
				if !client.IsErrConnectionFailed(err) {
					return err
				} else {
					util.Logger.Error(err)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

var loggingDrivers = map[string]struct{}{
	"local":     {},
	"json-file": {},
}

func isValidLoggingDriver(s string) bool {
	_, ok := loggingDrivers[s]
	return ok
}
