package xstring

import "testing"

func TestUUIDX(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UUIDX()
			t.Logf("UUIDX() = %v", got)
		})
	}
}
