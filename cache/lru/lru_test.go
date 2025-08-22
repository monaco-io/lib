package lru

import (
	"testing"
	"time"
)

var instance = New[int, int](100, time.Second*1)

func TestCache(t *testing.T) {
	instance.Set(1, 1)
	want, ok := instance.Get(1)
	if !ok {
		t.Fatal("want 1")
	}
	if want != 1 {
		t.Fatal("want 1")
	}
	time.Sleep(time.Second * 2)
	_, ok = instance.Get(1)
	if ok {
		t.Fatal("not want 1")
	}
	for i := 0; i < 110; i++ {
		instance.Set(i, i)
	}
	_, ok = instance.Get(1)
	if ok {
		t.Fatal("not want 1")
	}
	t.Log(instance.Len())
}
