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
	"deployment-manager/manager/itf"
	"github.com/docker/docker/api/types/mount"
)

var StateMap = map[string]itf.ContainerState{
	"created":    itf.UnknownState,
	"running":    itf.RunningState,
	"paused":     itf.UnknownState,
	"restarting": itf.UnknownState,
	"removing":   itf.UnknownState,
	"exited":     itf.StoppedState,
	"dead":       itf.UnhealthyState,
}

var StateRMap = map[itf.ContainerState]string{
	itf.RunningState:   "running",
	itf.StoppedState:   "exited",
	itf.UnhealthyState: "dead",
	itf.UnknownState:   "created",
}

var RestartPolicyMap = map[string]itf.RestartStrategy{
	"no":             itf.RestartNever,
	"on-failure":     itf.RestartOnFail,
	"always":         itf.RestartAlways,
	"unless-stopped": itf.RestartNotStopped,
}

var RestartPolicyRMap = func() map[itf.RestartStrategy]string {
	m := make(map[itf.RestartStrategy]string)
	for k, v := range RestartPolicyMap {
		m[v] = k
	}
	return m
}()

var MountTypeMap = map[mount.Type]itf.MountType{
	mount.TypeBind:   itf.BindMount,
	mount.TypeVolume: itf.VolumeMount,
	mount.TypeTmpfs:  itf.TmpfsMount,
}

var MountTypeRMap = func() map[itf.MountType]mount.Type {
	m := make(map[itf.MountType]mount.Type)
	for k, v := range MountTypeMap {
		m[v] = k
	}
	return m
}()

var PortTypeMap = map[string]itf.PortType{
	"tcp":  itf.TcpPort,
	"udp":  itf.UdpPort,
	"sctp": itf.SctpPort,
}

var PortTypeRMap = func() map[itf.PortType]string {
	m := make(map[itf.PortType]string)
	for k, v := range PortTypeMap {
		m[v] = k
	}
	return m
}()

var NetTypeMap = map[string]itf.NetworkType{
	"bridge":  itf.BridgeNet,
	"macvlan": itf.MACVlanNet,
	"host":    itf.HostNet,
}

var NetTypeRMap = func() map[itf.NetworkType]string {
	m := make(map[itf.NetworkType]string)
	for k, v := range NetTypeMap {
		m[v] = k
	}
	return m
}()
