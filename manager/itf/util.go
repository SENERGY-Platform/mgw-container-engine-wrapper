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

package itf

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func (p Port) KeyStr() string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

func (m Mount) KeyStr() string {
	return fmt.Sprintf("%s:%s", m.Source, m.Target)
}

func (s Subnet) KeyStr() string {
	return fmt.Sprintf("%s/%d", s.Prefix.String(), s.Bits)
}

func (p *PortType) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	if t, ok := PortTypeMap[s]; ok {
		*p = t
	} else {
		err = errors.New(fmt.Sprintf("unknown port type '%s'", s))
	}
	return
}

func (n *NetworkType) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	if t, ok := NetworkTypeMap[s]; ok {
		*n = t
	} else {
		err = errors.New(fmt.Sprintf("unknown network type '%s'", s))
	}
	return
}

func (r *RestartStrategy) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	if st, ok := RestartStrategyMap[s]; ok {
		*r = st
	} else {
		err = errors.New(fmt.Sprintf("unknown restart strategy '%s'", s))
	}
	return
}

func (m *MountType) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	if t, ok := MountTypeMap[s]; ok {
		*m = t
	} else {
		err = errors.New(fmt.Sprintf("unknown mount type '%s'", s))
	}
	return
}

func (i *IPAddr) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	if ip := net.ParseIP(s); ip != nil {
		i.IP = ip
	} else {
		err = errors.New(fmt.Sprintf("invalid IP address '%s'", s))
	}
	return
}
