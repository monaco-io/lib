package typing

import "testing"

var mp = NewSyncMap[int, int]()

func TestMap(t *testing.T) {
	t.Log(mp.Load(1))
	mp.Store(1, 1)
	t.Log(mp.Load(1))
}
