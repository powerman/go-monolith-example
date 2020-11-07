package def

import (
	"context"
	"sync"
	"time"
)

// A Cancellation is an interface capturing only the deadline and
// cancellation functionality of a context.
type Cancellation interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
}

type mergedContext struct {
	context.Context
	extra Cancellation
	done  <-chan struct{}
	sync.Mutex
	err error
}

// MergeCancel implements proposal https://github.com/golang/go/issues/36503.
func MergeCancel(parent context.Context, extra Cancellation) (ctx context.Context, cancel context.CancelFunc) {
	m := &mergedContext{
		Context: parent,
		extra:   extra,
	}

	cancelled := make(chan struct{}, 1)
	cancel = func() {
		select {
		case cancelled <- struct{}{}:
		default:
		}
	}

	done1 := parent.Done()
	done2 := extra.Done()
	switch {
	case done1 == nil:
		m.done = done2
	case done2 == nil:
		m.done = done1
	default:
		done := make(chan struct{})
		m.done = done
		go func() {
			select {
			case <-done1:
			case <-done2:
			case <-cancelled:
				m.Lock()
				defer m.Unlock()
				if m.err == nil {
					m.err = context.Canceled
				}
			}
			close(done)
		}()
	}

	return m, cancel
}

func (m *mergedContext) Deadline() (deadline time.Time, ok bool) {
	deadline1, ok1 := m.Context.Deadline()
	if deadline2, ok2 := m.extra.Deadline(); ok2 && (!ok1 || deadline2.Before(deadline1)) {
		return deadline2, ok2
	}
	return deadline1, ok1
}

func (m *mergedContext) Done() <-chan struct{} {
	return m.done
}

func (m *mergedContext) Err() error {
	m.Lock()
	defer m.Unlock()
	if m.err == nil {
		m.err = m.Context.Err()
	}
	if m.err == nil {
		m.err = m.extra.Err()
	}
	return m.err
}
