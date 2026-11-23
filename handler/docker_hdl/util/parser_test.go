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
	"github.com/docker/docker/api/types/container"
	"reflect"
	"strings"
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

func TestParseHealthOutput(t *testing.T) {
	str := func(s string) *string { return &s }
	res := func(out string) *container.HealthcheckResult { return &container.HealthcheckResult{Output: out} }
	long := strings.Repeat("a", 2000)
	// 'ä' is two bytes, so byte 1024 falls inside a rune when the prefix has an odd length
	multi := "x" + strings.Repeat("ä", 1000)
	tests := []struct {
		name string
		in   *container.Health
		want *string
	}{
		{name: "no health check", in: nil, want: nil},
		{name: "no log entry", in: &container.Health{}, want: nil},
		{name: "empty output", in: &container.Health{Log: []*container.HealthcheckResult{res("")}}, want: str("")},
		{name: "trimmed", in: &container.Health{Log: []*container.HealthcheckResult{res("  broken\n")}}, want: str("broken")},
		{name: "newest entry wins", in: &container.Health{Log: []*container.HealthcheckResult{res("old"), res("new\nsecond line")}}, want: str("new\nsecond line")},
		{name: "cut at limit", in: &container.Health{Log: []*container.HealthcheckResult{res(long)}}, want: str(long[:1024])},
		{name: "cut on rune boundary", in: &container.Health{Log: []*container.HealthcheckResult{res(multi)}}, want: str("x" + strings.Repeat("ä", 511))},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseHealthOutput(tc.in)
			if (got == nil) != (tc.want == nil) || (got != nil && *got != *tc.want) {
				t.Errorf("ParseHealthOutput() = %v, want %v", got, tc.want)
			}
		})
	}
}
