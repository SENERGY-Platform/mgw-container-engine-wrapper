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
	"github.com/SENERGY-Platform/mgw-deployment-manager-lib/dm-lib"
	"github.com/docker/docker/api/types/mount"
)

var StateMap = map[string]dm_lib.ContainerState{
	"created":    dm_lib.UnknownState,
	"running":    dm_lib.RunningState,
	"paused":     dm_lib.UnknownState,
	"restarting": dm_lib.UnknownState,
	"removing":   dm_lib.UnknownState,
	"exited":     dm_lib.StoppedState,
	"dead":       dm_lib.UnhealthyState,
}

var StateRMap = map[dm_lib.ContainerState]string{
	dm_lib.RunningState:   "running",
	dm_lib.StoppedState:   "exited",
	dm_lib.UnhealthyState: "dead",
	dm_lib.UnknownState:   "created",
}

var RestartPolicyMap = map[string]dm_lib.RestartStrategy{
	"no":             dm_lib.RestartNever,
	"on-failure":     dm_lib.RestartOnFail,
	"always":         dm_lib.RestartAlways,
	"unless-stopped": dm_lib.RestartNotStopped,
}

var RestartPolicyRMap = func() map[dm_lib.RestartStrategy]string {
	m := make(map[dm_lib.RestartStrategy]string)
	for k, v := range RestartPolicyMap {
		m[v] = k
	}
	return m
}()

var MountTypeMap = map[mount.Type]dm_lib.MountType{
	mount.TypeBind:   dm_lib.BindMount,
	mount.TypeVolume: dm_lib.VolumeMount,
	mount.TypeTmpfs:  dm_lib.TmpfsMount,
}

var MountTypeRMap = func() map[dm_lib.MountType]mount.Type {
	m := make(map[dm_lib.MountType]mount.Type)
	for k, v := range MountTypeMap {
		m[v] = k
	}
	return m
}()

var PortTypeMap = map[string]dm_lib.PortType{
	"tcp":  dm_lib.TcpPort,
	"udp":  dm_lib.UdpPort,
	"sctp": dm_lib.SctpPort,
}

var PortTypeRMap = func() map[dm_lib.PortType]string {
	m := make(map[dm_lib.PortType]string)
	for k, v := range PortTypeMap {
		m[v] = k
	}
	return m
}()

var NetTypeMap = map[string]dm_lib.NetworkType{
	"bridge":  dm_lib.BridgeNet,
	"macvlan": dm_lib.MACVlanNet,
	"host":    dm_lib.HostNet,
}

var NetTypeRMap = func() map[dm_lib.NetworkType]string {
	m := make(map[dm_lib.NetworkType]string)
	for k, v := range NetTypeMap {
		m[v] = k
	}
	return m
}()
