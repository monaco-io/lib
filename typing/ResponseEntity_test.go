package typing_test

import (
	"testing"

	"github.com/monaco-io/lib/typing"
)

func TestResponseEntity_ToString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "case 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r typing.ResponseEntity
			got := r.ToString()
			t.Logf("ToString() = %v", got)
		})
	}
}
