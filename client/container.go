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
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetContainers(ctx context.Context, filter model.ContainerFilter) ([]model.Container, error) {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath)
	if err != nil {
		return nil, err
	}
	u += genGetContainersQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var containers []model.Container
	err = c.baseClient.ExecRequestJSON(req, &containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c *Client) GetContainer(ctx context.Context, id string) (model.Container, error) {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath, id)
	if err != nil {
		return model.Container{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.Container{}, err
	}
	var container model.Container
	err = c.baseClient.ExecRequestJSON(req, &container)
	if err != nil {
		return model.Container{}, err
	}
	return container, nil
}

func (c *Client) CreateContainer(ctx context.Context, container model.Container) (id string, err error) {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(container)
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

func (c *Client) StartContainer(ctx context.Context, id string) error {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath, id, model.ContainerStartPath)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, nil)
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}

func (c *Client) StopContainer(ctx context.Context, id string) (jobId string, err error) {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath, id, model.ContainerStopPath)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, nil)
	if err != nil {
		return "", err
	}
	return c.baseClient.ExecRequestString(req)
}

func (c *Client) RestartContainer(ctx context.Context, id string) (jobId string, err error) {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath, id, model.ContainerRestartPath)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, nil)
	if err != nil {
		return "", err
	}
	return c.baseClient.ExecRequestString(req)
}

func (c *Client) RemoveContainer(ctx context.Context, id string, force bool) error {
	u, err := url.JoinPath(c.baseUrl, model.ContainersPath, id)
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

func (c *Client) GetContainerLog(ctx context.Context, id string, logOptions model.LogFilter) (io.ReadCloser, error) {
	panic("not implemented")
}

func genGetContainersQuery(filter model.ContainerFilter) string {
	var q []string
	if filter.Name != "" {
		q = append(q, "name="+filter.Name)
	}
	if filter.State != "" {
		q = append(q, "state="+filter.State)
	}
	if len(filter.Labels) > 0 {
		for k, v := range filter.Labels {
			if v != "" {
				q = append(q, "label="+k+"="+v)
			} else {
				q = append(q, "label="+k)
			}
		}
	}
	if len(q) > 0 {
		return "?" + strings.Join(q, "&")
	}
	return ""
}
