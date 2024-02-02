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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetVolumes(ctx context.Context, filter model.VolumeFilter) ([]model.Volume, error) {
	u, err := url.JoinPath(c.baseUrl, model.VolumesPath)
	if err != nil {
		return nil, err
	}
	u += genVolumesQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var volumes []model.Volume
	err = c.baseClient.ExecRequestJSON(req, &volumes)
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

func (c *Client) GetVolume(ctx context.Context, id string) (model.Volume, error) {
	u, err := url.JoinPath(c.baseUrl, model.VolumesPath, id)
	if err != nil {
		return model.Volume{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.Volume{}, err
	}
	var volume model.Volume
	err = c.baseClient.ExecRequestJSON(req, &volume)
	if err != nil {
		return model.Volume{}, err
	}
	return volume, nil
}

func (c *Client) CreateVolume(ctx context.Context, vol model.Volume) (string, error) {
	u, err := url.JoinPath(c.baseUrl, model.VolumesPath)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(vol)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return c.baseClient.ExecRequestString(req)
}

func (c *Client) RemoveVolume(ctx context.Context, id string, force bool) error {
	u, err := url.JoinPath(c.baseUrl, model.VolumesPath, id)
	if err != nil {
		return err
	}
	if force {
		u += "?force=true"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}

func genVolumesQuery(filter model.VolumeFilter) string {
	var q []string
	if len(filter.Labels) > 0 {
		q = append(q, "labels="+genLabels(filter.Labels, "=", ","))
	}
	if len(q) > 0 {
		return "?" + strings.Join(q, "&")
	}
	return ""
}
