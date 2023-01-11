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

const (
	TcpPort  PortType = "tcp"
	UdpPort  PortType = "udp"
	SctpPort PortType = "sctp"
)

var PortTypeMap = map[string]PortType{
	string(TcpPort):  TcpPort,
	string(UdpPort):  UdpPort,
	string(SctpPort): SctpPort,
}

const (
	BridgeNet  NetworkType = "bridge"
	MACVlanNet NetworkType = "macvlan"
	HostNet    NetworkType = "host"
)

var NetworkTypeMap = map[string]NetworkType{
	string(BridgeNet):  BridgeNet,
	string(MACVlanNet): MACVlanNet,
	string(HostNet):    HostNet,
}

const (
	RestartNever      RestartStrategy = "never"
	RestartAlways     RestartStrategy = "always"
	RestartNotStopped RestartStrategy = "not-stopped"
	RestartOnFail     RestartStrategy = "on-fail"
)

var RestartStrategyMap = map[string]RestartStrategy{
	string(RestartNever):      RestartNever,
	string(RestartAlways):     RestartAlways,
	string(RestartNotStopped): RestartNotStopped,
	string(RestartOnFail):     RestartOnFail,
}

const (
	BindMount   MountType = "bind"
	VolumeMount MountType = "volume"
	TmpfsMount  MountType = "tmpfs"
)

var MountTypeMap = map[string]MountType{
	string(BindMount):   BindMount,
	string(VolumeMount): VolumeMount,
	string(TmpfsMount):  TmpfsMount,
}

const (
	InitState       ContainerState = "initialized"
	RunningState    ContainerState = "running"
	RestartingState ContainerState = "restarting"
	StoppedState    ContainerState = "stopped"
	UnhealthyState  ContainerState = "unhealthy"
	UnknownState    ContainerState = "unknown"
)

var ContainerStateMap = map[string]ContainerState{
	string(InitState):       InitState,
	string(RunningState):    RunningState,
	string(RestartingState): RestartingState,
	string(StoppedState):    StoppedState,
	string(UnhealthyState):  UnhealthyState,
	string(UnknownState):    UnknownState,
}

const (
	ContainerStart   ContainerSetState = "start"
	ContainerStop    ContainerSetState = "stop"
	ContainerRestart ContainerSetState = "restart"
)
