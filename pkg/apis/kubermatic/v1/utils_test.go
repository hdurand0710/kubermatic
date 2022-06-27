/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package v1

import (
	"net"
	"reflect"
	"testing"
)

func expectedResult(start, end int) *NetworkRanges {
	res := NetworkRanges{}
	res.CIDRBlocks = make([]string, 0)
	for i := start; i <= end; i++ {
		ipnet := net.IPNet{
			IP:   net.IPv4(172, 10, byte(i), 0).To4(),
			Mask: net.CIDRMask(24, 32),
		}
		res.CIDRBlocks = append(res.CIDRBlocks, ipnet.String())
	}
	return &res
}

func TestSplitCIDRByNumber(t *testing.T) {
	type args struct {
		cidrBlocks []string
		number     uint32
	}
	tests := []struct {
		name    string
		args    args
		want    *NetworkRanges
		wantErr bool
	}{
		{
			name: "mask/16to/24",
			args: args{
				cidrBlocks: []string{"172.10.0.0/16"},
				// Split into 256 addresses per range to get /24 subnets.
				// Expected result is:
				// - 172.10.0.0/24
				// - 172.10.1.0/24
				// ...
				// - 172.10.255.0/24
				number: 256,
			},
			want:    expectedResult(0, 255),
			wantErr: false,
		},
		{
			name: "error-can-not-split",
			args: args{
				cidrBlocks: []string{"172.10.0.0/24"},
				// Split into 512 addresses per range to get /23 subnets.
				// Expected result an error, subranges are bigger that the CIDR, can not split.
				number: 512,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "mask/23to/24",
			args: args{
				cidrBlocks: []string{"172.10.2.0/23"},
				// Split into 256 addresses per range to get /24 subnets.
				// Expected result is:
				// - 172.10.2.0/24
				// - 172.10.3.0/24
				number: 256,
			},
			want:    expectedResult(2, 3),
			wantErr: false,
		},
		{
			name: "mask/23to/24 - 2 Blocks",
			args: args{
				cidrBlocks: []string{"172.10.2.0/23", "172.11.2.0/23"},
				// Split into 256 addresses per range to get /24 subnets.
				// Expected result is:
				// - 172.10.2.0/24
				// - 172.10.3.0/24
				// - 172.11.2.0/24
				// - 172.11.3.0/24
				number: 256,
			},
			want: &NetworkRanges{
				CIDRBlocks: []string{"172.10.2.0/24", "172.10.3.0/24", "172.11.2.0/24", "172.11.3.0/24"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NetworkRanges{
				CIDRBlocks: tt.args.cidrBlocks,
			}
			got, err := r.SplitByNumberOfAddresses(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("ß%s - error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s - res= %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
