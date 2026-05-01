package main

import (
	"sync"
	"testing"
)

func TestThreadSafe(t *testing.T) {
	counter := Counter{}

	var wg sync.WaitGroup
	wg.Add(1000)

	for range 1000 {
		go func() {
			counter.Inc()
			wg.Done()
		}()
	}

	wg.Wait()

	got := counter.Value()
	want := 1000

	if got != want {
		t.Errorf("got: %d\n want: %d\n", got, want)
	}

}
