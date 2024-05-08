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

func (c *Client) GetImages(ctx context.Context, filter model.ImageFilter) ([]model.Image, error) {
	u, err := url.JoinPath(c.baseUrl, model.ImagesPath)
	if err != nil {
		return nil, err
	}
	u += genImagesQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var images []model.Image
	err = c.baseClient.ExecRequestJSON(req, &images)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (c *Client) GetImage(ctx context.Context, id string) (model.Image, error) {
	u, err := url.JoinPath(c.baseUrl, model.ImagesPath, id)
	if err != nil {
		return model.Image{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.Image{}, err
	}
	var image model.Image
	err = c.baseClient.ExecRequestJSON(req, &image)
	if err != nil {
		return model.Image{}, err
	}
	return image, nil
}

func (c *Client) AddImage(ctx context.Context, img string) (jobId string, err error) {
	u, err := url.JoinPath(c.baseUrl, model.ImagesPath)
	if err != nil {
		return "", err
	}
	imgReq := model.ImageRequest{Image: img}
	body, err := json.Marshal(imgReq)
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

func (c *Client) RemoveImage(ctx context.Context, id string) error {
	u, err := url.JoinPath(c.baseUrl, model.ImagesPath, id)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}

func genImagesQuery(filter model.ImageFilter) string {
	var q []string
	if len(filter.Labels) > 0 {
		q = append(q, "labels="+genLabels(filter.Labels, "=", ","))
	}
	if filter.Name != "" {
		q = append(q, "name="+filter.Name)
	}
	if filter.Tag != "" {
		q = append(q, "tag="+filter.Tag)
	}
	if len(q) > 0 {
		return "?" + strings.Join(q, "&")
	}
	return ""
}
