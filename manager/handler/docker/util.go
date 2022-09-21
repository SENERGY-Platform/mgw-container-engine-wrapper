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
	"bytes"
	"deployment-manager/manager/itf"
	"deployment-manager/util"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"net"
	"strconv"
	"strings"
	"time"
)

func parseEndpointSettings(endptSettings map[string]*network.EndpointSettings) (netInfo []itf.ContainerNet) {
	if len(endptSettings) > 0 {
		for key, val := range endptSettings {
			netInfo = append(netInfo, itf.ContainerNet{
				ID:          val.NetworkID,
				Name:        key,
				IPAddress:   itf.IPAddr{IP: net.ParseIP(val.IPAddress)},
				Gateway:     itf.IPAddr{IP: net.ParseIP(val.Gateway)},
				DomainNames: val.Aliases,
				MacAddress:  val.MacAddress,
			})
		}
	}
	return
}

func parsePortSetAndMap(portSet nat.PortSet, portMap nat.PortMap) (ports []itf.Port) {
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
					Interface: itf.IPAddr{IP: net.ParseIP(binding.HostIP)},
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

func parseMountPoints(mountPoints []types.MountPoint) (mounts []itf.Mount) {
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

func parseMounts(mts []mount.Mount) (mounts []itf.Mount) {
	if len(mts) > 0 {
		for _, mt := range mts {
			if mType, ok := mountTypeMap[mt.Type]; ok {
				m := itf.Mount{
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
					m.Mode = mt.TmpfsOptions.Mode
				}
				mounts = append(mounts, m)
			}
		}
	}
	return
}

func parseEnv(ev []string) (env map[string]string) {
	if len(ev) > 0 {
		env = make(map[string]string, len(ev))
		for _, s := range ev {
			p := strings.Split(s, "=")
			env[p[0]] = p[1]
		}
	}
	return
}

func genEnv(ev map[string]string) (env []string) {
	if len(ev) > 0 {
		for key, val := range ev {
			env = append(env, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return
}

func parseStopTimeout(t *int) *itf.Duration {
	if t != nil {
		return &itf.Duration{Duration: time.Duration(*t * int(time.Second))}
	}
	return nil
}

func genStopTimeout(d *itf.Duration) *int {
	if d != nil {
		t := int(d.Seconds())
		return &t
	}
	return nil
}

func genPortMap(ports []itf.Port) (nat.PortMap, error) {
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
				HostIP:   binding.Interface.String(),
				HostPort: strconv.FormatInt(int64(binding.Number), 10),
			})
		}
		pm[port] = bindings
	}
	return pm, nil
}

func genMounts(mounts []itf.Mount) ([]mount.Mount, error) {
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

func genFilterArgs(filter [][2]string) (f filters.Args) {
	if filter != nil && len(filter) > 0 {
		f = filters.NewArgs()
		for _, i := range filter {
			f.Add(i[0], i[1])
		}
	}
	return
}

func parseNetIPAMConfig(c []network.IPAMConfig) (s itf.Subnet, gw itf.IPAddr) {
	if c != nil && len(c) > 0 {
		sp := strings.Split(c[0].Subnet, "/")
		if len(sp) == 2 {
			s.Prefix.IP = net.ParseIP(sp[0])
			i, _ := strconv.ParseInt(sp[1], 10, 0)
			s.Bits = int(i)
		}
		gw.IP = net.ParseIP(c[0].Gateway)
	}
	return
}

func genNetIPAMConfig(n itf.Network) (c []network.IPAMConfig) {
	c = append(c, network.IPAMConfig{
		Subnet:  n.Subnet.KeyStr(),
		Gateway: n.Gateway.String(),
	})
	return
}

func parseTimestamp(s string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, s)
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

type ImgPullResp struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	ID             string `json:"id"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func (r ImgPullResp) String() string {
	var b bytes.Buffer
	if r.Status != "" {
		b.WriteString(r.Status)
	}
	if r.ID != "" {
		b.WriteString(" " + r.ID)
	}
	if r.Message != "" {
		b.WriteString(" " + r.Message)
	}
	if strings.Contains(r.Status, "Downloading") {
		b.WriteString(fmt.Sprintf(" %d/%d", r.ProgressDetail.Current, r.ProgressDetail.Total))
	}
	return b.String()
}
