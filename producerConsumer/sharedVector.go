package main

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type sharedVector struct {
	mutex       sync.Mutex
	empty, full *semaphore.Weighted
	vector      []int
	index       int
	inserted    int
	historic    []int
}

// newSharedVector creates a new shared vector with n elements.
func newSharedVector(n int) *sharedVector {
	fullSemaphore := semaphore.NewWeighted(int64(n - 1))

	for i := 0; i < n-1; i++ {
		fullSemaphore.Acquire(context.Background(), 1)
	}

	return &sharedVector{
		vector:   make([]int, n),
		empty:    semaphore.NewWeighted(int64(n - 1)),
		full:     fullSemaphore,
		historic: make([]int, 0, 2*MAX_CONSUMED+50),
	}
}

// insert adds a value in the shared vector.
func (sv *sharedVector) insert(value int) bool {
	if sv.inserted > MAX_CONSUMED {
		return false
	}

	sv.empty.Acquire(context.Background(), 1)
	sv.mutex.Lock()
	sv.index++

	sv.vector[sv.index] = value
	sv.inserted++

	sv.historic = append(sv.historic, sv.index)
	sv.mutex.Unlock()
	sv.full.Release(1)

	return true
}

// pop removes a value from the shared vector
func (sv *sharedVector) pop() (int, bool) {
	if sv.inserted > MAX_CONSUMED && sv.index < 1 {
		return 0, false
	}

	sv.full.Acquire(context.Background(), 1)
	sv.mutex.Lock()
	value := sv.vector[sv.index]
	sv.index--

	sv.historic = append(sv.historic, sv.index)
	sv.mutex.Unlock()
	sv.empty.Release(1)

	return value, true
}
