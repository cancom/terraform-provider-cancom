package util

import "context"

// lock can be used to lock, but make it `context.Context` aware.
// e.g. it'll respect cancelling and timeouts.
type lock chan struct{}

func NewLock() lock {
	return make(lock, 1)

}

func (c lock) Lock(ctx context.Context) error {
	select {
	case c <- struct{}{}:
		// lock acquired
		return nil
	case <-ctx.Done():
		// Timeout
		return ctx.Err()
	}
}

func (c lock) Unlock() {
	<-c
}
