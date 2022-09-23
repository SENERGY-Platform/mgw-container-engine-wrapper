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
	"fmt"
	"net"
	"time"
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
		err = fmt.Errorf("unknown port type '%s'", s)
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
		err = fmt.Errorf("unknown network type '%s'", s)
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
		err = fmt.Errorf("unknown restart strategy '%s'", s)
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
		err = fmt.Errorf("unknown mount type '%s'", s)
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
		d.Duration = dur
	}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func NewError(code uint, msg string, err error) error {
	return &Error{
		code: code,
		msg:  msg,
		err:  err,
	}
}

func (e *Error) Error() string {
	if e.msg != "" && e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else if e.msg != "" {
		return e.msg
	} else {
		return e.err.Error()
	}
}

func (e *Error) Unwrap() error {
	return e.err
}
