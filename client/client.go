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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"io"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HttpClient
	baseUrl    string
}

func New(httpClient HttpClient, baseUrl string) *Client {
	return &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func execRequest(httpClient HttpClient, req *http.Request) ([]byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		msg := resp.Status
		if len(body) > 0 {
			msg = string(body)
		}
		return nil, getError(resp.StatusCode, msg)
	}
	return body, nil
}

func execRequestJSONResp(httpClient HttpClient, req *http.Request, v any) error {
	body, err := execRequest(httpClient, req)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func getError(sc int, msg string) error {
	err := errors.New(msg)
	switch sc {
	case http.StatusInternalServerError:
		return model.NewInternalError(err)
	case http.StatusNotFound:
		return model.NewNotFoundError(err)
	case http.StatusBadRequest:
		return model.NewInvalidInputError(err)
	}
	return newResponseError(sc, err)
}
