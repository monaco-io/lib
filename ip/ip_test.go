package ip

import (
	"reflect"
	"testing"
)

func TestInternalV4(t *testing.T) {
	tests := []struct {
		name   string
		wantIp string
	}{
		{
			name:   "bilibili-local",
			wantIp: "10.23.2.138",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIp := InternalV4(); gotIp != tt.wantIp {
				t.Errorf("Internal() = %v, want %v", gotIp, tt.wantIp)
			}
		})
	}
}

func TestExternalV4(t *testing.T) {
	tests := []struct {
		name    string
		wantIps []string
	}{
		{
			name:    "bilibili-out",
			wantIps: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIps := ExternalV4(); !reflect.DeepEqual(gotIps, tt.wantIps) {
				t.Errorf("External() = %v, want %v", gotIps, tt.wantIps)
			}
		})
	}
}

func TestInetAtoN(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
	}{
		{
			args: args{
				s: "10.23.2.138",
			},
			wantSum: 169280138,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSum := AtoI(tt.args.s); gotSum != tt.wantSum {
				t.Errorf("InetAtoN() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestInetNtoA(t *testing.T) {
	type args struct {
		sum int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				sum: 169280138,
			},
			want: "10.23.2.138",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ItoA(tt.args.sum); got != tt.want {
				t.Errorf("InetNtoA() = %v, want %v", got, tt.want)
			}
		})
	}
}
