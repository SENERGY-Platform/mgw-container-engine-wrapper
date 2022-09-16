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

package docker

import (
	"deployment-manager/manager/ce-handler/itf"
	"deployment-manager/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"strconv"
	"strings"
)

func parseContainerNetworks(endptSettings map[string]*network.EndpointSettings) (netInfo []itf.ContainerNet) {
	if len(endptSettings) > 0 {
		for key, val := range endptSettings {
			netInfo = append(netInfo, itf.ContainerNet{
				ID:          val.NetworkID,
				Name:        key,
				IPAddress:   val.IPAddress,
				Gateway:     val.Gateway,
				DomainNames: val.Aliases,
				MacAddress:  val.MacAddress,
			})
		}
	}
	return
}

func parseContainerPorts(portSet nat.PortSet, portMap nat.PortMap) (ports []itf.Port) {
	if len(portSet) > 0 || len(portMap) > 0 {
		set := make(map[string]struct{})
		for port, bindings := range portMap {
			p := itf.Port{
				Number:   port.Int(),
				Protocol: portTypeMap[port.Proto()],
			}
			for _, binding := range bindings {
				num, err := strconv.ParseInt(binding.HostPort, 10, 0)
				if err != nil {
					util.Logger.Error(err)
				}
				p.Bindings = append(p.Bindings, itf.PortBinding{
					Number:    int(num),
					Interface: binding.HostIP,
				})
			}
			ports = append(ports, p)
			set[p.String()] = struct{}{}
		}
		for port, _ := range portSet {
			if _, ok := set[string(port)]; !ok {
				ports = append(ports, itf.Port{
					Number:   port.Int(),
					Protocol: portTypeMap[port.Proto()],
				})
			}
		}
	}
	return
}

func parseContainerMounts(mountPoints []types.MountPoint) (mounts []itf.Mount) {
	if len(mountPoints) > 0 {
		for _, mp := range mountPoints {
			if mType, ok := mountTypeMap[mp.Type]; ok {
				mounts = append(mounts, itf.Mount{
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

func parseContainerEnvVars(ev []string) (env map[string]string) {
	if len(ev) > 0 {
		env = make(map[string]string, len(ev))
		for _, s := range ev {
			p := strings.Split(s, "=")
			env[p[0]] = p[1]
		}
	}
	return
}

func parseRestartPolicy(rp container.RestartPolicy) itf.RestartConfig {
	return itf.RestartConfig{
		Strategy: restartPolicyMap[rp.Name],
		Retries:  rp.MaximumRetryCount,
	}
}
