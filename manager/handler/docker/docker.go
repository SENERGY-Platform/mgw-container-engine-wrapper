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
	"github.com/docker/docker/client"
)

type Docker struct {
	client *client.Client
}

func New(ops ...client.Opt) (*Docker, error) {
	c, err := client.NewClientWithOpts(ops...)
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
