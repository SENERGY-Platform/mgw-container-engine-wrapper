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
	"context"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) GetJobs(ctx context.Context, filter model.JobFilter) ([]model.Job, error) {
	u, err := url.JoinPath(c.baseUrl, model.JobsPath)
	if err != nil {
		return nil, err
	}
	u += genJobsFilter(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var jobs []model.Job
	err = execRequestJSONResp(c.httpClient, req, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (c *Client) GetJob(ctx context.Context, id string) (model.Job, error) {
	panic("not implemented")
}

func (c *Client) CancelJob(ctx context.Context, id string) error {
	panic("not implemented")
}

func genJobsFilter(filter model.JobFilter) string {
	var q []string
	if filter.SortDesc {
		q = append(q, "sort_desc=true")
	}
	if filter.Status != "" {
		q = append(q, "status="+filter.Status)
	}
	if !filter.Since.IsZero() {
		q = append(q, "since="+filter.Since.Format(time.RFC3339Nano))
	}
	if !filter.Until.IsZero() {
		q = append(q, "until="+filter.Until.Format(time.RFC3339Nano))
	}
	if len(q) > 0 {
		return "?" + strings.Join(q, "&")
	}
	return ""
}
