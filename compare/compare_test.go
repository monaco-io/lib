package compare

import (
	"testing"
)

func TestMax(t *testing.T) {
	t.Log(Max(1, 2) == 2)
	t.Log(Max(-9999, 2, 888, 77, 1, 0, 9999) == 9999)
}

func TestMin(t *testing.T) {
	t.Log(Min(1, 2) == 1)
	t.Log(Min(-9999, 2, 888, 77, 1, 0) == -9999)
}

func TestContains(t *testing.T) {
	t.Log(Contains(1, 0, 1, 3, 4))
	var a, b *int
	_a := 1
	_b := 2
	a, b = &_a, &_b
	t.Log(Contains(&_a, a, b))
}
