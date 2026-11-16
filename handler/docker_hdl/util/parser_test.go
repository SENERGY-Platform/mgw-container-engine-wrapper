/*
 * Copyright 2026 InfAI (CC SES)
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
	"reflect"
	"testing"
)

func TestParseEnv(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want map[string]string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "key and value", in: []string{"A=1", "B=two"}, want: map[string]string{"A": "1", "B": "two"}},
		{name: "empty value", in: []string{"A="}, want: map[string]string{"A": ""}},
		{name: "entry without equals sign", in: []string{"REACT_APP_API_URL", "A=1"}, want: map[string]string{"REACT_APP_API_URL": "", "A": "1"}},
		{name: "value containing equals signs", in: []string{"OPTS=a=1,b=2"}, want: map[string]string{"OPTS": "a=1,b=2"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ParseEnv(tc.in); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseEnv(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
