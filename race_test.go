package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

/*
* race condition using mutex
* cons : with multiple states, very complicated and slower
***/
func TestRaceConditionWithMutex(t *testing.T) {
	var state int32
	var mu sync.RWMutex

	for i := 0; i < 10; i++ {
		go func(i int) {
			mu.Lock()
			state += int32(i)
			mu.Unlock()
		}(i)
	}
}

/*
* race condition using atomic values
* we could use atomic.Value, but this example the logic of updating is atomic not the state itslef
***/
func TestRaceConditionWithAtomicValues(t *testing.T) {
	var state int32

	for i := 0; i < 10; i++ {
		go func(i int) {
			atomic.AddInt32(&state, int32(i))
		}(i)
	}
}
