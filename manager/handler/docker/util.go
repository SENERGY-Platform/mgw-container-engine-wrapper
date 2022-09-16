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
	"deployment-manager/manager/itf"
	"deployment-manager/util"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"strconv"
	"strings"
	"time"
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
			set[p.KeyStr()] = struct{}{}
		}
		for port := range portSet {
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

func genContainerEnvVars(ev map[string]string) (env []string) {
	if len(ev) > 0 {
		for key, val := range ev {
			env = append(env, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return
}

func parseStopTimeout(t *int) *time.Duration {
	if t != nil {
		d := time.Duration(*t * int(time.Second))
		return &d
	}
	return nil
}

func getStopTimeout(d *time.Duration) *int {
	if d != nil {
		t := int(d.Seconds())
		return &t
	}
	return nil
}

func getPorts(ports []itf.Port) (nat.PortMap, error) {
	pm := make(nat.PortMap)
	set := make(map[string]struct{})
	for _, p := range ports {
		if _, ok := set[p.KeyStr()]; ok {
			return pm, errors.New("port duplicate")
		}
		set[p.KeyStr()] = struct{}{}
		port, err := nat.NewPort(portTypeRMap[p.Protocol], strconv.FormatInt(int64(p.Number), 10))
		if err != nil {
			return pm, err
		}
		var bindings []nat.PortBinding
		for _, binding := range p.Bindings {
			bindings = append(bindings, nat.PortBinding{
				HostIP:   binding.Interface,
				HostPort: strconv.FormatInt(int64(binding.Number), 10),
			})
		}
		pm[port] = bindings
	}
	return pm, nil
}

func getMounts(mounts []itf.Mount) ([]mount.Mount, error) {
	var msl []mount.Mount
	set := make(map[string]struct{})
	for i := 0; i < len(mounts); i++ {
		m := mounts[i]
		if _, ok := set[m.KeyStr()]; ok {
			return msl, errors.New("mount duplicate")
		}
		set[m.KeyStr()] = struct{}{}
		mnt := mount.Mount{
			Type:     mountTypeRMap[m.Type],
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		}
		switch m.Type {
		case itf.VolumeMount:
			mnt.VolumeOptions = &mount.VolumeOptions{Labels: m.Labels}
		case itf.TmpfsMount:
			mnt.TmpfsOptions = &mount.TmpfsOptions{
				SizeBytes: m.Size,
				Mode:      m.Mode,
			}
		}
		msl = append(msl, mnt)
	}
	return msl, nil
}

func checkNetworks(n []itf.ContainerNet) error {
	set := make(map[string]struct{})
	for _, net := range n {
		if _, ok := set[net.Name]; ok {
			return errors.New("network duplicate")
		}
		set[net.Name] = struct{}{}
	}
	return nil
}
