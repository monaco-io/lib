package xec

import "testing"

func TestRegister(t *testing.T) {
	type args struct {
		code int
		msg  map[Language]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				code: 0,
				msg: map[Language]string{
					LangDefault: "success",
					LangEnUS:    "成功",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.code, tt.args.msg)
			t.Log(_messages.Load().(Message))
		})
	}
}
