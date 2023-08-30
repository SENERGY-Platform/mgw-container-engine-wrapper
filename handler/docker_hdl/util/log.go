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

type RCWrapper struct {
	ReadCloser     io.ReadCloser
	notMultiplexed bool
	remainder      int
}

func (c *RCWrapper) Close() error {
	return c.ReadCloser.Close()
}

func (c *RCWrapper) Read(p []byte) (n int, err error) {
	pLen := len(p)
	if !c.notMultiplexed {
		var pRemainder int
		for n < pLen {
			pRemainder = pLen - n
			var size int
			var n2 int
			if c.remainder == 0 {
				header := make([]byte, 8)
				if n2, err = c.ReadCloser.Read(header); err != nil {
					break
				}
				streamType := header[0]
				if streamType < 1 || streamType > 2 {
					if !c.notMultiplexed {
						c.notMultiplexed = true
						n = n2
						copy(p, header)
						break
					}
					err = fmt.Errorf("unkown stream type '%d'", streamType)
					break
				}
				outputLen := int(binary.BigEndian.Uint32(header[4:]))
				if outputLen < pRemainder {
					size = outputLen
				} else {
					size = pRemainder
					c.remainder = outputLen - pRemainder
				}
			} else {
				if c.remainder < pRemainder {
					size = c.remainder
					c.remainder = 0
				} else {
					size = pRemainder
					c.remainder = c.remainder - pRemainder
				}
			}
			n2, err = c.ReadCloser.Read(p[n : n+size])
			n += n2
			if err != nil {
				break
			}
		}
	}
	if c.notMultiplexed {
		for n < pLen {
			var n2 int
			n2, err = c.ReadCloser.Read(p[n:pLen])
			n += n2
			if err != nil {
				break
			}
		}
	}
	return
}
