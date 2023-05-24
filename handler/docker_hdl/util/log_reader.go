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
	"encoding/binary"
	"fmt"
	"io"
)

type LogReader struct {
	rc     io.ReadCloser
	remain int
}

func NewLogReader(rc io.ReadCloser) *LogReader {
	return &LogReader{
		rc: rc,
	}
}

func (r *LogReader) Read(p []byte) (n int, err error) {
	pSize := len(p)
	var pRemain int
	for n < pSize {
		pRemain = pSize - n
		var size int
		var iType byte
		if r.remain == 0 {
			var head = make([]byte, 8)
			if _, err = r.rc.Read(head); err != nil {
				break
			}
			iType = head[0]
			if iType < 1 || iType > 2 {
				err = fmt.Errorf("unkown input type '%d'", iType)
				break
			}
			oSize := int(binary.BigEndian.Uint32(head[4:]))
			if oSize < pRemain {
				size = oSize
			} else {
				size = pRemain
				r.remain = oSize - pRemain
			}
		} else {
			if r.remain < pRemain {
				size = r.remain
				r.remain = 0
			} else {
				size = pRemain
				r.remain = r.remain - pRemain
			}
		}
		var n2 int
		n2, err = r.rc.Read(p[n : n+size])
		n += n2
		if err != nil {
			break
		}
	}
	return
}

func (r *LogReader) Close() error {
	return r.rc.Close()
}
