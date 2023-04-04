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
	"container-engine-wrapper/model"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"net"
	"os"
	"strconv"
	"time"
)

func GenEnv(ev map[string]string) (env []string) {
	if len(ev) > 0 {
		for key, val := range ev {
			env = append(env, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return
}

func GenStopTimeout(d *model.Duration) *int {
	if d != nil {
		t := int(time.Duration(*d).Seconds())
		return &t
	}
	return nil
}

func GenPortMap(ports []model.Port) (nat.PortMap, error) {
	pm := make(nat.PortMap)
	set := make(map[string]struct{})
	for _, p := range ports {
		if _, ok := model.PortTypeMap[p.Protocol]; !ok {
			return pm, fmt.Errorf("invalid port type '%s'", p.Protocol)
		}
		if _, ok := set[p.KeyStr()]; ok {
			return pm, fmt.Errorf("port duplicate '%s'", p.KeyStr())
		}
		set[p.KeyStr()] = struct{}{}
		port, err := nat.NewPort(PortTypeRMap[p.Protocol], strconv.FormatInt(int64(p.Number), 10))
		if err != nil {
			return pm, err
		}
		var bindings []nat.PortBinding
		for _, binding := range p.Bindings {
			bindings = append(bindings, nat.PortBinding{
				HostIP:   net.IP(binding.Interface).String(),
				HostPort: strconv.FormatInt(int64(binding.Number), 10),
			})
		}
		pm[port] = bindings
	}
	return pm, nil
}

func GenMounts(mounts []model.Mount) ([]mount.Mount, error) {
	var msl []mount.Mount
	set := make(map[string]struct{})
	for i := 0; i < len(mounts); i++ {
		m := mounts[i]
		if _, ok := model.MountTypeMap[m.Type]; !ok {
			return msl, fmt.Errorf("invalid mount type '%s'", m.Type)
		}
		if _, ok := set[m.KeyStr()]; ok {
			return msl, fmt.Errorf("mount duplicate '%s'", m.KeyStr())
		}
		set[m.KeyStr()] = struct{}{}
		mnt := mount.Mount{
			Type:     MountTypeRMap[m.Type],
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		}
		switch m.Type {
		case model.VolumeMount:
			mnt.VolumeOptions = &mount.VolumeOptions{Labels: m.Labels}
		case model.TmpfsMount:
			mnt.TmpfsOptions = &mount.TmpfsOptions{
				SizeBytes: m.Size,
				Mode:      os.FileMode(m.Mode),
			}
		}
		msl = append(msl, mnt)
	}
	return msl, nil
}

func genLabelFilterArgs(fArgs *filters.Args, fLabels map[string]string) {
	if fLabels != nil && len(fLabels) > 0 {
		for k, v := range fLabels {
			l := k
			if v != "" {
				l += "=" + v
			}
			fArgs.Add("label", l)
		}
	}
}

func GenContainerFilterArgs(filter model.ContainerFilter) filters.Args {
	fArgs := filters.NewArgs()
	if filter.Name != "" {
		fArgs.Add("name", filter.Name)
	}
	if filter.State != "" {
		fArgs.Add("status", StateRMap[filter.State])
	}
	genLabelFilterArgs(&fArgs, filter.Labels)
	return fArgs
}

func GenImageFilterArgs(filter model.ImageFilter) filters.Args {
	fArgs := filters.NewArgs()
	genLabelFilterArgs(&fArgs, filter.Labels)
	return fArgs
}

func GenVolumeFilterArgs(filter model.VolumeFilter) filters.Args {
	fArgs := filters.NewArgs()
	genLabelFilterArgs(&fArgs, filter.Labels)
	return fArgs
}

func GenNetIPAMConfig(n model.Network) (c []network.IPAMConfig) {
	c = append(c, network.IPAMConfig{
		Subnet:  n.Subnet.KeyStr(),
		Gateway: net.IP(n.Gateway).String(),
	})
	return
}

func GenRestartPolicy(strategy model.RestartStrategy, retries *int) (rp container.RestartPolicy, err error) {
	if _, ok := model.RestartStrategyMap[strategy]; !ok {
		err = fmt.Errorf("invalid restart strategy '%s'", strategy)
		return
	}
	if strategy == model.RestartOnFail && retries == nil {
		err = fmt.Errorf("invalid restart strategy configuration: number of retries = %v", retries)
		return
	}
	rp.Name = RestartPolicyRMap[strategy]
	if retries != nil {
		rp.MaximumRetryCount = *retries
	}
	return
}

func CheckNetworks(n []model.ContainerNet) error {
	set := make(map[string]struct{})
	for _, net := range n {
		if _, ok := set[net.Name]; ok {
			return fmt.Errorf("network duplicate '%s'", net.Name)
		}
		set[net.Name] = struct{}{}
	}
	return nil
}
