package retry

import (
	"context"
	"errors"
	"testing"
)

func TestDo(t *testing.T) {
	type args struct {
		f    func(context.Context) error
		opts Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				f: func(context.Context) error { panic("not implemented") },
				opts: Options{
					RetryTimes: 3,
					Context:    nil,
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				f: func(context.Context) error { return errors.New("mock error") },
				opts: Options{
					RetryTimes: 10,
					Context:    nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Do(tt.args.f, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
