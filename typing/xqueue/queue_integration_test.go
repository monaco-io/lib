package xqueue

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueueIntegration(t *testing.T) {
	var processed int64

	// 创建队列
	q := New(func(data int) error {
		atomic.AddInt64(&processed, 1)
		time.Sleep(time.Millisecond) // 模拟处理时间
		return nil
	})

	// 发送数据
	const numItems = 100
	for i := 0; i < numItems; i++ {
		q.Input(i)
	}

	// 等待一段时间让数据处理
	time.Sleep(time.Millisecond * 100)

	// 关闭队列
	q.CloseSync()

	// 验证所有数据都被处理了
	processedCount := atomic.LoadInt64(&processed)
	if processedCount != numItems {
		t.Errorf("Expected %d items to be processed, got %d", numItems, processedCount)
	}
}

func TestQueueConcurrentInputAndClose(t *testing.T) {
	var processed int64

	q := New(func(data int) error {
		atomic.AddInt64(&processed, 1)
		return nil
	})

	// 并发发送数据
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				q.Input(start*10 + j)
			}
		}(i)
	}

	// 等待发送完成
	wg.Wait()

	// 给一些时间让队列处理数据
	time.Sleep(time.Millisecond * 50)

	// 关闭队列
	q.CloseSync()

	// 验证处理的数据数量
	processedCount := atomic.LoadInt64(&processed)
	if processedCount != 100 {
		t.Errorf("Expected 100 items to be processed, got %d", processedCount)
	}
}
