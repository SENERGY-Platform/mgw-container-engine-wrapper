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

const ServiceName = "mgw-container-engine-wrapper"

const (
	TcpPort  PortType = "tcp"
	UdpPort  PortType = "udp"
	SctpPort PortType = "sctp"
)

var PortTypeMap = map[PortType]struct{}{
	TcpPort:  {},
	UdpPort:  {},
	SctpPort: {},
}

const (
	BridgeNet  NetworkType = "bridge"
	MACVlanNet NetworkType = "macvlan"
	HostNet    NetworkType = "host"
)

var NetworkTypeMap = map[NetworkType]struct{}{
	BridgeNet:  {},
	MACVlanNet: {},
	HostNet:    {},
}

const (
	RestartNever      RestartStrategy = "never"
	RestartAlways     RestartStrategy = "always"
	RestartNotStopped RestartStrategy = "not-stopped"
	RestartOnFail     RestartStrategy = "on-fail"
)

var RestartStrategyMap = map[RestartStrategy]struct{}{
	RestartNever:      {},
	RestartAlways:     {},
	RestartNotStopped: {},
	RestartOnFail:     {},
}

const (
	BindMount   MountType = "bind"
	VolumeMount MountType = "volume"
	TmpfsMount  MountType = "tmpfs"
)

var MountTypeMap = map[MountType]struct{}{
	BindMount:   {},
	VolumeMount: {},
	TmpfsMount:  {},
}

const (
	InitState       ContainerState = "initialized"
	RunningState    ContainerState = "running"
	PausedState     ContainerState = "paused"
	RestartingState ContainerState = "restarting"
	RemovingState   ContainerState = "removing"
	StoppedState    ContainerState = "stopped"
	UnhealthyState  ContainerState = "unhealthy"
)

var ContainerStateMap = map[ContainerState]struct{}{
	InitState:       {},
	RunningState:    {},
	PausedState:     {},
	RestartingState: {},
	RemovingState:   {},
	StoppedState:    {},
	UnhealthyState:  {},
}

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCanceled  JobStatus = "canceled"
	JobCompleted JobStatus = "completed"
	JobError     JobStatus = "error"
	JobOK        JobStatus = "ok"
)

var JobStateMap = map[JobStatus]struct{}{
	JobPending:   {},
	JobRunning:   {},
	JobCanceled:  {},
	JobCompleted: {},
	JobError:     {},
	JobOK:        {},
}

const (
	ContainersPath    = "containers"
	ContainerCtrlPath = "ctrl"
	ContainerLogsPath = "logs"
	ImagesPath        = "images"
	NetworksPath      = "networks"
	VolumesPath       = "volumes"
	JobsPath          = "jobs"
	JobsCancelPath    = "cancel"
)
