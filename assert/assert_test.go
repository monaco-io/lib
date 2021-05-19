package assert

import (
	"testing"
)

func TestOK(t *testing.T) {
	const a = 1
	type args struct {
		ok  bool
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		// {
		// 	args: args{
		// 		ok:  1 == 2,
		// 		msg: "want false",
		// 	},
		// },
		{
			args: args{
				ok:  a >= 1,
				msg: "want true",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OK(tt.args.ok, tt.args.msg)
		})
	}
}

func TestDeepEqual(t *testing.T) {
	type args struct {
		x interface{}
		y interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				x: nil,
				y: nil,
			},
			want: true,
		},
		{
			args: args{
				x: 1,
				y: 1,
			},
			want: true,
		},
		{
			args: args{
				x: 1,
				y: "1",
			},
			want: false,
		},
		{
			args: args{
				x: args{},
				y: args{},
			},
			want: true,
		},
		{
			args: args{
				x: args{
					x: 1,
				},
				y: args{
					x: 1,
				},
			},
			want: true,
		},
		{
			args: args{
				x: map[string]string{
					"x": "y",
				},
				y: map[string]string{
					"x": "y",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeepEqual(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("DeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		x interface{}
		y interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				x: nil,
				y: nil,
			},
			want: true,
		},
		{
			args: args{
				x: 1,
				y: 1,
			},
			want: true,
		},
		{
			args: args{
				x: 1,
				y: "1",
			},
			want: false,
		},
		{
			args: args{
				x: args{
					x: 1,
				},
				y: args{
					x: 1,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
