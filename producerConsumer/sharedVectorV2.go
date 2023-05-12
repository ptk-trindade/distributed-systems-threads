package main

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// V2
type sharedVectorV2 struct {
	mtxConsumer sync.Mutex
	mtxProducer sync.Mutex
	empty, full *semaphore.Weighted
	vector      []int
	cosumeIndex int
	insertIndex int

	ctx       context.Context
	cancelCtx context.CancelFunc

	insertTime []time.Time
	popTime    []time.Time
}

func newSharedVectorV2(n int) *sharedVectorV2 {
	ctx, cancelCtx := context.WithCancel(context.Background())

	fullSemaphore := semaphore.NewWeighted(int64(n))

	fullSemaphore.Acquire(ctx, int64(n))

	return &sharedVectorV2{
		vector:     make([]int, n),
		empty:      semaphore.NewWeighted(int64(n)),
		full:       fullSemaphore,
		insertTime: make([]time.Time, 0, MAX_CONSUMED+50),
		popTime:    make([]time.Time, 0, MAX_CONSUMED+50),
		ctx:        ctx,
		cancelCtx:  cancelCtx,
	}
}

func (sv *sharedVectorV2) insert(value int) bool {
	length := len(sv.vector)

	err := sv.empty.Acquire(sv.ctx, 1)
	if err != nil {
		return false
	}
	sv.mtxProducer.Lock()
	if sv.insertIndex >= MAX_CONSUMED {
		sv.cancelCtx()
		sv.mtxProducer.Unlock()
		return false
	}

	sv.vector[sv.insertIndex%length] = value
	sv.insertIndex++

	sv.insertTime = append(sv.insertTime, time.Now())

	sv.mtxProducer.Unlock()
	sv.full.Release(1)

	return true
}

func (sv *sharedVectorV2) pop() (int, bool) {
	length := len(sv.vector)

	err := sv.full.Acquire(sv.ctx, 1)
	if err != nil {
		return 0, false
	}
	sv.mtxConsumer.Lock()

	value := sv.vector[sv.cosumeIndex%length]
	sv.cosumeIndex++

	sv.popTime = append(sv.popTime, time.Now())

	sv.mtxConsumer.Unlock()
	sv.empty.Release(1)

	return value, true
}

func (sv *sharedVectorV2) getHistoric() []int {
	historic := make([]int, 1, len(sv.insertTime)+len(sv.popTime)+1)
	var i, j int

	for i < len(sv.insertTime) || j < len(sv.popTime) {
		if i < len(sv.insertTime) && j < len(sv.popTime) {
			if sv.popTime[j].Before(sv.insertTime[i]) { // pop
				historic = append(historic, historic[len(historic)-1]-1)
				j++
			} else if sv.insertTime[i].Before(sv.popTime[j]) { // insert
				historic = append(historic, historic[len(historic)-1]+1)
				i++
			} else { // insert == pop
				historic = append(historic, historic[len(historic)-1])
				i++
				j++
			}

		} else if i < len(sv.insertTime) { // insert
			historic = append(historic, historic[len(historic)-1]+1)
			i++
		} else { // pop
			historic = append(historic, historic[len(historic)-1]-1)
			j++
		}
	}

	return historic
}
