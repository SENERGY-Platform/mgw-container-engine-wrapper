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

const (
	TcpPort  PortType = "tcp"
	UdpPort  PortType = "udp"
	SctpPort PortType = "sctp"
)

const (
	BridgeNet  NetworkType = "bridge"
	MACVlanNet NetworkType = "macvlan"
	HostNet    NetworkType = "host"
)

const (
	RunningState   ContainerState = "running"
	StoppedState   ContainerState = "stopped"
	UnhealthyState ContainerState = "unhealthy"
	UnknownState   ContainerState = "unknown"
)

const (
	RestartNever      RestartStrategy = "never"
	RestartAlways     RestartStrategy = "always"
	RestartNotStopped RestartStrategy = "not-stopped"
	RestartOnFail     RestartStrategy = "on-fail"
)

const (
	BindMount   MountType = "bind"
	VolumeMount MountType = "volume"
	TmpfsMount  MountType = "tmpfs"
)
