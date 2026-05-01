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

	want := 1000

	AssertCounter(t, &counter, want)

}

func AssertCounter(t *testing.T, got *Counter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got: %d\n want: %d\n", got.Value(), want)
	}
}
