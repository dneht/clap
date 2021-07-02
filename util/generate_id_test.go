package util

import (
	"sync"
	"testing"
)

func init() {
	testing.Init()
}

const size = 1000000
const thread = 64

var sm sync.Map

func TestUniqueId(t *testing.T) {
	channel := make(chan bool, thread)
	defer close(channel)

	wg := sync.WaitGroup{}
	for i := 0; i < size; i++ {
		wg.Add(1)
		channel <- true
		go func() {
			id := UniqueId()
			sm.Store(id, true)
			wg.Done()
			<-channel
		}()
	}

	wg.Wait()
	smSize := 0
	sm.Range(func(key, value interface{}) bool {
		smSize++
		return true
	})
	if smSize != size {
		t.Error("check error", smSize)
	}
}
