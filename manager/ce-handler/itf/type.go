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

package itf

import (
	"context"
)

type ContainerEngineHandler interface {
	ListContainers(ctx context.Context, filter map[string]string) (map[string]*Container, error)
	ContainerInfo(ctx context.Context, id string) (*Container, error)
	ImageInfo(ctx context.Context, id string) (*Image, error)
}

type Image struct {
	ID      string   `json:"id"`
	Created string   `json:"created"`
	Size    int64    `json:"size"`
	Arch    string   `json:"arch"`
	Tags    []string `json:"tags"`
	Digests []string `json:"digests"`
}

type Network struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Port struct {
	Number   int    `json:"number"`
	Protocol string `json:"protocol"`
	Bindings []PortBinding
}

type PortBinding struct {
	Number    int    `json:"number"`
	Interface string `json:"interface"`
}

type NetworkInfo struct {
	NetworkID   string   `json:"network_id"`
	IPAddress   string   `json:"ip_address"`
	Gateway     string   `json:"gateway"`
	DomainNames []string `json:"domain_names"`
	MacAddress  string   `json:"mac_address"`
}

type Mount struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Target string `json:"target"`
	Mode   string `json:"mode"`
}

type RestartStrategy int

type RestartConfig struct {
	Strategy RestartStrategy `json:"strategy"`
	Retries  int             `json:"retries"`
}

type Container struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	ImageID       string                 `json:"image_id"`
	Image         string                 `json:"image"`
	State         string                 `json:"state"`
	RestartConfig RestartConfig          `json:"restart_config"`
	Created       string                 `json:"created"`
	Started       string                 `json:"started"`
	Hostname      string                 `json:"hostname"`
	Env           map[string]string      `json:"env"`
	Mounts        map[string]Mount       `json:"mounts"`
	Labels        map[string]string      `json:"labels"`
	Ports         map[string]Port        `json:"ports"`
	Networks      map[string]NetworkInfo `json:"networks"`
}
