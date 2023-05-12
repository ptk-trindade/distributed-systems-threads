package main

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type sharedVector interface {
	insert(value int) bool
	pop() (int, bool)
	getHistoric() []int
	// ctx 		context.Context
	// cancelCtx	context.CancelFunc
}

// V1
type sharedVectorV1 struct {
	mutex       sync.Mutex
	empty, full *semaphore.Weighted
	vector      []int
	index       int
	inserted    int
	historic    []int
	ctx         context.Context
	cancelCtx   context.CancelFunc
}

// newSharedVector creates a new shared vector with n elements.
func newSharedVectorV1(n int) *sharedVectorV1 {
	ctx, cancelCtx := context.WithCancel(context.Background())
	fullSemaphore := semaphore.NewWeighted(int64(n))

	fullSemaphore.Acquire(ctx, int64(n))

	return &sharedVectorV1{
		vector:    make([]int, n),
		empty:     semaphore.NewWeighted(int64(n)),
		full:      fullSemaphore,
		historic:  make([]int, 0, 2*MAX_CONSUMED+50),
		ctx:       ctx,
		cancelCtx: cancelCtx,
	}
}

// insert adds a value in the shared vector.
func (sv *sharedVectorV1) insert(value int) bool {

	err := sv.empty.Acquire(sv.ctx, 1)
	if err != nil {
		return false
	}

	keepGoing := true
	sv.mutex.Lock()
	if sv.inserted >= MAX_CONSUMED {
		sv.cancelCtx()
		sv.mutex.Unlock()
		return false
	}

	sv.vector[sv.index] = value
	sv.index++
	sv.inserted++

	sv.historic = append(sv.historic, sv.index)
	sv.mutex.Unlock()
	sv.full.Release(1)

	return keepGoing
}

// pop removes a value from the shared vector
func (sv *sharedVectorV1) pop() (int, bool) {
	err := sv.full.Acquire(sv.ctx, 1) // returns err if context is Done
	if err != nil {
		return 0, false
	}

	sv.mutex.Lock()
	sv.index--
	value := sv.vector[sv.index]

	sv.historic = append(sv.historic, sv.index)
	sv.mutex.Unlock()
	sv.empty.Release(1)

	return value, true
}

func (sv *sharedVectorV1) getHistoric() []int {
	return sv.historic
}
