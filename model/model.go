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

package model

import (
	"io/fs"
	"net"
	"time"
)

type Duration time.Duration

type FileMode fs.FileMode

type Time time.Time

type IPAddr net.IP

type NetworkType = string

type PortType = string

type MountType = string

type RestartStrategy = string

type ContainerState = string

// Image ---------------------------------------------------------------------------------------

type Image struct {
	ID      string            `json:"id"`
	Created Time              `json:"created"`
	Size    int64             `json:"size"`
	Arch    string            `json:"arch"`
	Tags    []string          `json:"tags"`
	Digests []string          `json:"digests"`
	Labels  map[string]string `json:"labels"`
}

type ImagesPostRequest struct {
	Image string `json:"image"`
}

// Network -------------------------------------------------------------------------------------

type Subnet struct {
	Prefix IPAddr `json:"prefix"`
	Bits   int    `json:"bits"`
}

type Network struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Type    NetworkType `json:"type"`
	Subnet  Subnet      `json:"subnet"`
	Gateway IPAddr      `json:"gateway"`
}

type Port struct {
	Number   int           `json:"number"`
	Protocol PortType      `json:"protocol"`
	Bindings []PortBinding `json:"bindings"`
}

type PortBinding struct {
	Number    int    `json:"number"`
	Interface IPAddr `json:"interface"`
}

// Container -----------------------------------------------------------------------------------

type Mount struct {
	Type     MountType         `json:"type"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	ReadOnly bool              `json:"read_only"`
	Labels   map[string]string `json:"labels,omitempty"`
	Size     int64             `json:"size,omitempty"`
	Mode     FileMode          `json:"mode,omitempty"`
}

type RunConfig struct {
	RestartStrategy RestartStrategy `json:"restart_strategy"`
	Retries         *int            `json:"retries"`
	RemoveAfterRun  bool            `json:"remove_after_run"`
	StopTimeout     *Duration       `json:"stop_timeout"`
	StopSignal      *string         `json:"stop_signal"`
	PseudoTTY       bool            `json:"pseudo_tty"`
}

type Container struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	State     ContainerState    `json:"state"`
	Created   Time              `json:"created"`
	Started   *Time             `json:"started"`
	Image     string            `json:"image"`
	ImageID   string            `json:"image_id"`
	EnvVars   map[string]string `json:"env_vars"`
	Labels    map[string]string `json:"labels"`
	Mounts    []Mount           `json:"mounts"`
	Ports     []Port            `json:"ports"`
	Networks  []ContainerNet    `json:"networks"`
	RunConfig RunConfig         `json:"run_config"`
}

type ContainerNet struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DomainNames []string `json:"domain_names"`
	Gateway     IPAddr   `json:"gateway"`
	IPAddress   IPAddr   `json:"ip_address"`
	MacAddress  string   `json:"mac_address"`
}

type ContainersPostResponse struct {
	ID string `json:"id"`
}

// Volume --------------------------------------------------------------------------------------

type Volume struct {
	Name    string            `json:"name"`
	Created Time              `json:"created"`
	Labels  map[string]string `json:"labels"`
}

// Job -----------------------------------------------------------------------------------------

type Job struct {
	ID        string        `json:"id"`
	Request   JobOrgRequest `json:"request"`
	Error     any           `json:"error"`
	Created   Time          `json:"created"`
	Started   *Time         `json:"started"`
	Completed *Time         `json:"completed"`
	Canceled  *Time         `json:"canceled"`
}

type JobOrgRequest struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
	Body   any    `json:"body"`
}
