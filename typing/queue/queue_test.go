package queue

import (
	"errors"
	"testing"

	"github.com/monaco-io/lib/typing/option"
)

func TestNew(t *testing.T) {
	type T struct {
		Value int
	}
	type args struct {
		consumer func(data T) error
		opts     []option.Option[Config]
	}
	tests := []struct {
		name string
		args args
		want Queue[T]
	}{
		{
			name: "create new queue with consumer",
			args: args{
				consumer: func(data T) error {
					t.Logf("consumer: %v", data)
					if data.Value == 50 {
						return errors.New("err=50")
					}
					return nil
				},
				opts: []option.Option[Config]{},
			},
			want: nil, // This would need to be updated based on actual Queue implementation
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.consumer, tt.args.opts...)
			for i := range 100 {
				got.Input(T{Value: i})
			}
			got.CloseSync()
		})
	}
}
