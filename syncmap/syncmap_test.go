package syncmap

import "testing"

var mp = New[int, int]()

func TestMap(t *testing.T) {
	t.Log(mp.Load(1))
	mp.Store(1, 1)
	t.Log(mp.Load(1))
}
