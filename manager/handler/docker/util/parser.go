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
	"github.com/SENERGY-Platform/mgw-container-engine-manager/manager/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"net"
	"strconv"
	"strings"
	"time"
)

func ParseEndpointSettings(endptSettings map[string]*network.EndpointSettings) (netInfo []model.ContainerNet) {
	if len(endptSettings) > 0 {
		for key, val := range endptSettings {
			netInfo = append(netInfo, model.ContainerNet{
				ID:          val.NetworkID,
				Name:        key,
				IPAddress:   model.IPAddr(net.ParseIP(val.IPAddress)),
				Gateway:     model.IPAddr(net.ParseIP(val.Gateway)),
				DomainNames: val.Aliases,
				MacAddress:  val.MacAddress,
			})
		}
	}
	return
}

func ParsePortSetAndMap(portSet nat.PortSet, portMap nat.PortMap) ([]model.Port, error) {
	var ports []model.Port
	if len(portSet) > 0 || len(portMap) > 0 {
		set := make(map[string]struct{})
		for port, bindings := range portMap {
			p := model.Port{
				Number:   port.Int(),
				Protocol: PortTypeMap[port.Proto()],
			}
			for _, binding := range bindings {
				num, err := strconv.ParseInt(binding.HostPort, 10, 0)
				if err != nil {
					return ports, err
				}
				p.Bindings = append(p.Bindings, model.PortBinding{
					Number:    int(num),
					Interface: model.IPAddr(net.ParseIP(binding.HostIP)),
				})
			}
			ports = append(ports, p)
			set[p.KeyStr()] = struct{}{}
		}
		for port := range portSet {
			if _, ok := set[string(port)]; !ok {
				ports = append(ports, model.Port{
					Number:   port.Int(),
					Protocol: PortTypeMap[port.Proto()],
				})
			}
		}
	}
	return ports, nil
}

func ParseMountPoints(mountPoints []types.MountPoint) (mounts []model.Mount) {
	if len(mountPoints) > 0 {
		for _, mp := range mountPoints {
			if mType, ok := MountTypeMap[mp.Type]; ok {
				mounts = append(mounts, model.Mount{
					Type:     mType,
					Source:   mp.Source,
					Target:   mp.Destination,
					ReadOnly: !mp.RW,
				})
			}
		}
	}
	return
}

func ParseMounts(mts []mount.Mount) (mounts []model.Mount) {
	if len(mts) > 0 {
		for _, mt := range mts {
			if mType, ok := MountTypeMap[mt.Type]; ok {
				m := model.Mount{
					Type:     mType,
					Source:   mt.Source,
					Target:   mt.Target,
					ReadOnly: mt.ReadOnly,
				}
				if mt.VolumeOptions != nil {
					m.Labels = mt.VolumeOptions.Labels
				}
				if mt.TmpfsOptions != nil {
					m.Size = mt.TmpfsOptions.SizeBytes
					m.Mode = model.FileMode(mt.TmpfsOptions.Mode)
				}
				mounts = append(mounts, m)
			}
		}
	}
	return
}

func ParseEnv(ev []string) (env map[string]string) {
	if len(ev) > 0 {
		env = make(map[string]string, len(ev))
		for _, s := range ev {
			p := strings.Split(s, "=")
			env[p[0]] = p[1]
		}
	}
	return
}

func ParseStopTimeout(t *int) *model.Duration {
	if t != nil {
		d := model.Duration(time.Duration(*t * int(time.Second)))
		return &d
	}
	return nil
}

func ParseNetIPAMConfig(c []network.IPAMConfig) (s model.Subnet, gw model.IPAddr) {
	if c != nil && len(c) > 0 {
		sp := strings.Split(c[0].Subnet, "/")
		if len(sp) == 2 {
			s.Prefix = model.IPAddr(net.ParseIP(sp[0]))
			i, _ := strconv.ParseInt(sp[1], 10, 0)
			s.Bits = int(i)
		}
		gw = model.IPAddr(net.ParseIP(c[0].Gateway))
	}
	return
}

func ParseTimestamp(s string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, s)
}

func ParseContainerName(s string) string {
	return strings.TrimPrefix(s, "/")
}

func ParseRestartPolicy(rp container.RestartPolicy) (strategy model.RestartStrategy, retires *int) {
	if rp.Name == "" {
		strategy = model.RestartNever
	} else {
		strategy = RestartPolicyMap[rp.Name]
	}
	if strategy == model.RestartOnFail {
		retires = &rp.MaximumRetryCount
	}
	return
}
