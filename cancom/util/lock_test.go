package util

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLock_MutualExclusion(t *testing.T) {
	l := NewLock()
	var active int32
	var maxActive int32
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Lock(context.Background()); err != nil {
				t.Errorf("Lock failed: %v", err)
				return
			}
			cur := atomic.AddInt32(&active, 1)
			for {
				prev := atomic.LoadInt32(&maxActive)
				if cur <= prev || atomic.CompareAndSwapInt32(&maxActive, prev, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt32(&active, -1)
			l.Unlock()
		}()
	}

	wg.Wait()
	if maxActive != 1 {
		t.Fatalf("expected max concurrent holders to be 1, got %d", maxActive)
	}
}

func TestLock_ContextTimeout(t *testing.T) {
	l := NewLock()
	if err := l.Lock(context.Background()); err != nil {
		t.Fatalf("initial Lock failed: %v", err)
	}
	defer l.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := l.Lock(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}
