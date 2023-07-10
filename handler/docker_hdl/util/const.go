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

package util

import (
	"github.com/SENERGY-Platform/mgw-container-engine-wrapper/lib/model"
	"github.com/docker/docker/api/types/mount"
)

var StateMap = map[string]model.ContainerState{
	"created":    model.InitState,
	"running":    model.RunningState,
	"paused":     model.PausedState,
	"restarting": model.RestartingState,
	"removing":   model.RemovingState,
	"exited":     model.StoppedState,
	"dead":       model.DeadState,
}

var StateRMap = func() map[model.ContainerState]string {
	m := make(map[model.ContainerState]string)
	for k, v := range StateMap {
		m[v] = k
	}
	return m
}()

var RestartPolicyMap = map[string]model.RestartStrategy{
	"no":             model.RestartNever,
	"on-failure":     model.RestartOnFail,
	"always":         model.RestartAlways,
	"unless-stopped": model.RestartNotStopped,
}

var RestartPolicyRMap = func() map[model.RestartStrategy]string {
	m := make(map[model.RestartStrategy]string)
	for k, v := range RestartPolicyMap {
		m[v] = k
	}
	return m
}()

var MountTypeMap = map[mount.Type]model.MountType{
	mount.TypeBind:   model.BindMount,
	mount.TypeVolume: model.VolumeMount,
	mount.TypeTmpfs:  model.TmpfsMount,
}

var MountTypeRMap = func() map[model.MountType]mount.Type {
	m := make(map[model.MountType]mount.Type)
	for k, v := range MountTypeMap {
		m[v] = k
	}
	return m
}()

var PortTypeMap = map[string]model.PortType{
	"tcp":  model.TcpPort,
	"udp":  model.UdpPort,
	"sctp": model.SctpPort,
}

var PortTypeRMap = func() map[model.PortType]string {
	m := make(map[model.PortType]string)
	for k, v := range PortTypeMap {
		m[v] = k
	}
	return m
}()

var NetTypeMap = map[string]model.NetworkType{
	"bridge":  model.BridgeNet,
	"macvlan": model.MACVlanNet,
	"host":    model.HostNet,
}

var NetTypeRMap = func() map[model.NetworkType]string {
	m := make(map[model.NetworkType]string)
	for k, v := range NetTypeMap {
		m[v] = k
	}
	return m
}()

func GetConst(s string, m map[string]string) string {
	if c, ok := m[s]; ok {
		return c
	} else {
		return "unknown"
	}
}
