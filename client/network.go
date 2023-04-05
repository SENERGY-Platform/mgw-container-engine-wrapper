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
)

func (c *Client) GetNetworks(ctx context.Context) ([]model.Network, error) {
	u, err := url.JoinPath(c.baseUrl, model.NetworksPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var networks []model.Network
	err = execRequestJSONResp(c.httpClient, req, &networks)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

func (c *Client) GetNetwork(ctx context.Context, id string) (model.Network, error) {
	u, err := url.JoinPath(c.baseUrl, model.NetworksPath, id)
	if err != nil {
		return model.Network{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.Network{}, err
	}
	var network model.Network
	err = execRequestJSONResp(c.httpClient, req, &network)
	if err != nil {
		return model.Network{}, err
	}
	return network, nil
}

func (c *Client) CreateNetwork(ctx context.Context, net model.Network) (string, error) {
	u, err := url.JoinPath(c.baseUrl, model.NetworksPath)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(net)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	body, err = execRequest(c.httpClient, req)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) RemoveNetwork(ctx context.Context, id string) error {
	u, err := url.JoinPath(c.baseUrl, model.NetworksPath, id)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	_, err = execRequest(c.httpClient, req)
	if err != nil {
		return err
	}
	return nil
}
