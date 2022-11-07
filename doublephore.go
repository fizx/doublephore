package doublephore

import (
	"context"
	"errors"

	"golang.org/x/sync/semaphore"
)

type Doublephore struct {
	primary *semaphore.Weighted
	waiters *semaphore.Weighted
}

func NewWeighted(n, waitSize int64) *Doublephore {
	return &Doublephore{semaphore.NewWeighted(n), semaphore.NewWeighted(n + waitSize)}
}

func (s *Doublephore) Acquire(ctx context.Context, n int64) error {
	ok := s.waiters.TryAcquire(n)
	if !ok {
		return errors.New("too many waiters")
	}
	return s.primary.Acquire(ctx, n)
}

func (s *Doublephore) TryAcquire(n int64) bool {
	ok := s.waiters.TryAcquire(n)
	if !ok {
		return false
	}
	return s.primary.TryAcquire(n)
}

func (s *Doublephore) Release(n int64) {
	s.primary.Release(n)
	s.waiters.Release(n)
}
