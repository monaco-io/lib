package sys

import "testing"

func TestRecover(t *testing.T) {
	defer Recover()
	panic("test")
}
