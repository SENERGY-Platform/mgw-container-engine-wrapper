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
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib"
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"io"
	"net/http"
)

type CewClient = lib.Api

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

func (c *Client) execRequest(req *http.Request) (io.ReadCloser, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		errMsg, err := readString(resp.Body)
		if err != nil || errMsg == "" {
			errMsg = resp.Status
		}
		return nil, getError(resp.StatusCode, errMsg)
	}
	return resp.Body, nil
}

func (c *Client) execRequestJSON(req *http.Request, v any) error {
	body, err := c.execRequest(req)
	if err != nil {
		return err
	}
	defer body.Close()
	err = readJSON(body, v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) execRequestString(req *http.Request) (string, error) {
	body, err := c.execRequest(req)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return readString(body)
}

func (c *Client) execRequestVoid(req *http.Request) error {
	body, err := c.execRequest(req)
	if err != nil {
		return err
	}
	defer body.Close()
	_ = readVoid(body)
	return nil
}

func readVoid(rc io.ReadCloser) error {
	_, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	return nil
}

func readString(rc io.ReadCloser) (string, error) {
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readJSON(rc io.ReadCloser, v any) error {
	jd := json.NewDecoder(rc)
	err := jd.Decode(v)
	if err != nil {
		_ = readVoid(rc)
		return err
	}
	return nil
}

func getError(sc int, msg string) error {
	var err error
	err = newResponseError(sc, errors.New(msg))
	if sc < 500 {
		err = newClientError(err)
	}
	if sc >= 500 {
		err = newServerError(err)
	}
	switch sc {
	case http.StatusInternalServerError:
		err = model.NewInternalError(err)
	case http.StatusNotFound:
		err = model.NewNotFoundError(err)
	case http.StatusBadRequest:
		err = model.NewInvalidInputError(err)
	}
	return err
}
