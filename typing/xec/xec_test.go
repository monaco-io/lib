package xec

import (
	"errors"
	"testing"
)

func TestError_New(t *testing.T) {
	type fields struct {
		Message []string
		Code    int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"case 1", fields{[]string{"not found"}, 404}, "[404] [not found]"},
		{"case 2", fields{[]string{"internal error"}, 500}, "[500] [internal error]"},
		{"case 3", fields{[]string{"internal1", "internal2"}, 500}, "[500] [internal1 internal2]"},
		{"case 4", fields{[]string{}, 404}, "[404] []"},
		{"case 5", fields{nil, 404}, "[404] []"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.fields.Code, tt.fields.Message...)
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	type fields struct {
		Message []string
		Code    int
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "nil error",
			fields: fields{[]string{"not found"}, 404},
			args:   args{nil},
			want:   false,
		},
		{
			name:   "same error code",
			fields: fields{[]string{"not found"}, 404},
			args:   args{New(404, "resource not found")},
			want:   true,
		},
		{
			name:   "different error code",
			fields: fields{[]string{"not found"}, 404},
			args:   args{New(500, "internal error")},
			want:   false,
		},
		{
			name:   "non-Error type",
			fields: fields{[]string{"not found"}, 404},
			args:   args{errors.New("standard error")},
			want:   false,
		},
		{
			name:   "same error instance",
			fields: fields{[]string{"test1"}, 200},
			args:   args{New(200, "test")},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.fields.Code, tt.fields.Message...)
			if got := e.Is(tt.args.err); got != tt.want {
				t.Errorf("Error.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
