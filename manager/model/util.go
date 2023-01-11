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

package model

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"strconv"
	"time"
)

func (p *Port) KeyStr() string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

func (m *Mount) KeyStr() string {
	return fmt.Sprintf("%s:%s", m.Source, m.Target)
}

func (s *Subnet) KeyStr() string {
	return fmt.Sprintf("%s/%d", net.IP(s.Prefix).String(), s.Bits)
}

func (i *IPAddr) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	if ip := net.ParseIP(s); ip != nil {
		*i = IPAddr(ip)
	} else {
		err = fmt.Errorf("invalid IP address '%s'", s)
	}
	return
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if dur, err := time.ParseDuration(s); err != nil {
		return err
	} else {
		*d = Duration(dur)
	}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (m *FileMode) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	i, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return err
	}
	*m = FileMode(fs.FileMode(i))
	return nil
}

func (m FileMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(m), 8))
}
