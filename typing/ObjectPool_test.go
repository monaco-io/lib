package typing

import (
	"sync"
	"testing"
)

type testObject struct {
	value int
}

func TestObjectPool(t *testing.T) {
	// Create a new object pool for testObject
	pool := NewObjectPool[testObject]()

	// Test Get returns a non-nil object
	obj := pool.Get()
	if obj == nil {
		t.Fatal("Expected non-nil object from Get()")
	}

	// Set a value on the object
	obj.value = 42

	// Put the object back in the pool
	pool.Put(obj)

	// Get the object again and verify if it's reused
	obj2 := pool.Get()

	// Note: sync.Pool doesn't guarantee that we get the same object back
	// so we can only test that we get a valid object, not its value
	if obj2 == nil {
		t.Fatal("Expected non-nil object from second Get()")
	}
}

func TestObjectPoolConcurrent(t *testing.T) {
	pool := NewObjectPool[testObject]()
	const goroutines = 100
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Launch multiple goroutines to test concurrent access
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				// Get an object from the pool
				obj := pool.Get()
				if obj == nil {
					t.Errorf("goroutine %d: Got nil object at iteration %d", id, j)
					return
				}

				// Set a value
				obj.value = id*iterations + j

				// Put it back
				pool.Put(obj)
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkObjectPool(b *testing.B) {
	pool := NewObjectPool[testObject]()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := pool.Get()
			obj.value = 42
			pool.Put(obj)
		}
	})
}

// Benchmark comparing with and without using object pool
func BenchmarkWithWithoutPool(b *testing.B) {
	b.Run("WithPool", func(b *testing.B) {
		pool := NewObjectPool[testObject]()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			obj := pool.Get()
			obj.value = i
			pool.Put(obj)
		}
	})

	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obj := new(testObject)
			obj.value = i
			_ = obj
		}
	})
}
