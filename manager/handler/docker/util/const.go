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

package util

import (
	"github.com/SENERGY-Platform/mgw-container-engine-manager-lib/cem-lib"
	"github.com/docker/docker/api/types/mount"
)

var StateMap = map[string]cem_lib.ContainerState{
	"created":    cem_lib.InitState,
	"running":    cem_lib.RunningState,
	"paused":     cem_lib.UnknownState,
	"restarting": cem_lib.RestartingState,
	"removing":   cem_lib.UnknownState,
	"exited":     cem_lib.StoppedState,
	"dead":       cem_lib.UnhealthyState,
}

var StateRMap = func() map[cem_lib.ContainerState]string {
	m := make(map[cem_lib.ContainerState]string)
	for k, v := range StateMap {
		if v != cem_lib.UnknownState {
			m[v] = k
		}
	}
	return m
}()

var RestartPolicyMap = map[string]cem_lib.RestartStrategy{
	"no":             cem_lib.RestartNever,
	"on-failure":     cem_lib.RestartOnFail,
	"always":         cem_lib.RestartAlways,
	"unless-stopped": cem_lib.RestartNotStopped,
}

var RestartPolicyRMap = func() map[cem_lib.RestartStrategy]string {
	m := make(map[cem_lib.RestartStrategy]string)
	for k, v := range RestartPolicyMap {
		m[v] = k
	}
	return m
}()

var MountTypeMap = map[mount.Type]cem_lib.MountType{
	mount.TypeBind:   cem_lib.BindMount,
	mount.TypeVolume: cem_lib.VolumeMount,
	mount.TypeTmpfs:  cem_lib.TmpfsMount,
}

var MountTypeRMap = func() map[cem_lib.MountType]mount.Type {
	m := make(map[cem_lib.MountType]mount.Type)
	for k, v := range MountTypeMap {
		m[v] = k
	}
	return m
}()

var PortTypeMap = map[string]cem_lib.PortType{
	"tcp":  cem_lib.TcpPort,
	"udp":  cem_lib.UdpPort,
	"sctp": cem_lib.SctpPort,
}

var PortTypeRMap = func() map[cem_lib.PortType]string {
	m := make(map[cem_lib.PortType]string)
	for k, v := range PortTypeMap {
		m[v] = k
	}
	return m
}()

var NetTypeMap = map[string]cem_lib.NetworkType{
	"bridge":  cem_lib.BridgeNet,
	"macvlan": cem_lib.MACVlanNet,
	"host":    cem_lib.HostNet,
}

var NetTypeRMap = func() map[cem_lib.NetworkType]string {
	m := make(map[cem_lib.NetworkType]string)
	for k, v := range NetTypeMap {
		m[v] = k
	}
	return m
}()
