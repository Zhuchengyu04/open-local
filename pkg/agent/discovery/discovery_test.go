/*
Copyright © 2021 Alibaba Group Holding Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package discovery

import (
	"testing"
)

func TestFilterInfo(t *testing.T) {
	type args struct {
		in      []string
		include []string
		exclude []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				in:      []string{"paas1", "paas"},
				include: []string{"paas[0-9]*"},
				exclude: []string{"paas[0-1]+"},
			},
			want: []string{"paas"},
		},
		{
			name: "test2",
			args: args{
				in:      []string{"paas4", "share"},
				include: []string{"paas[0-9]+", "share"},
				exclude: []string{},
			},
			want: []string{"paas4", "share"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterInfo(tt.args.in, tt.args.include, tt.args.exclude); !sameStringSlice(got, tt.want) {
				t.Errorf("FilterInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
