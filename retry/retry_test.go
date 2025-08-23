package retry

import (
	"errors"
	"testing"

	"github.com/monaco-io/lib/typing/xopt"
)

func TestDo(t *testing.T) {
	type args struct {
		f    func() error
		opts func() xopt.Option[config]
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				f:    func() error { panic("not implemented") },
				opts: func() xopt.Option[config] { return WithRetryTimes(3) },
			},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				f:    func() error { return errors.New("mock error") },
				opts: func() xopt.Option[config] { return WithRetryTimes(10) },
			},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				f: func() error { return nil },
				opts: func() xopt.Option[config] {
					return WithRetryTimes(999)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Do(tt.args.f, tt.args.opts()); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
