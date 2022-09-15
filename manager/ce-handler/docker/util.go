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

func parseContainerNetworks(endptSettings map[string]*network.EndpointSettings) (netInfo map[string]itf.NetworkInfo) {
	if len(endptSettings) > 0 {
		netInfo = make(map[string]itf.NetworkInfo, len(endptSettings))
		for _, val := range endptSettings {
			netInfo[val.NetworkID] = itf.NetworkInfo{
				NetworkID:   val.NetworkID,
				IPAddress:   val.IPAddress,
				Gateway:     val.Gateway,
				DomainNames: val.Aliases,
				MacAddress:  val.MacAddress,
			}
		}
	}
	return
}

func parseContainerPorts(portSet nat.PortSet, portMap nat.PortMap) (ports map[string]itf.Port) {
	if len(portSet) > 0 || len(portMap) > 0 {
		ports = make(map[string]itf.Port, len(portSet))
		for port, bindings := range portMap {
			p := itf.Port{
				Number:   port.Int(),
				Protocol: port.Proto(),
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
			ports[p.String()] = p
		}
		for port, _ := range portSet {
			if _, ok := ports[string(port)]; !ok {
				p := itf.Port{
					Number:   port.Int(),
					Protocol: port.Proto(),
				}
				ports[p.String()] = p
			}
		}
	}
	return
}

func parseContainerMounts(mountPoints []types.MountPoint) (mounts map[string]itf.Mount) {
	if len(mountPoints) > 0 {
		mounts = make(map[string]itf.Mount, len(mountPoints))
		for _, mp := range mountPoints {
			m := itf.Mount{
				Type:     string(mp.Type),
				Source:   mp.Source,
				Target:   mp.Destination,
				ReadOnly: !mp.RW,
			}
			mounts[m.Target] = m
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
