package main

import (
	"fmt"
	"sync"
)

type H2O struct {
	// TODO: add synchronization primitives
	// e.g., mutex, condition variable, channels, semaphores (via buffered channels)
	cond    sync.Cond
	counter int
}

type Locker struct {
	mu sync.Mutex
}

func (l *Locker) Lock() {
	l.mu.Lock()
}

func (l *Locker) Unlock() {
	l.mu.Unlock()
}

func NewH2O() *H2O {
	return &H2O{
		*sync.NewCond(&Locker{sync.Mutex{}}),
		0,
	}
}

func (h2o *H2O) Hydrogen(releaseHydrogen func()) {
	h2o.cond.L.Lock()

	for h2o.counter >= 2 {
		h2o.cond.Wait()
	}

	releaseHydrogen()

	h2o.counter++

	h2o.cond.L.Unlock()
	h2o.cond.Broadcast()
}

func (h2o *H2O) Oxygen(releaseOxygen func()) {
	h2o.cond.L.Lock()

	for h2o.counter < 2 {
		h2o.cond.Wait()
	}

	releaseOxygen()

	h2o.counter++
	if h2o.counter >= 2 {
		h2o.counter = 0
	}

	h2o.cond.L.Unlock()
	h2o.cond.Broadcast()
}

func main() {
	h2o := NewH2O()
	var wg sync.WaitGroup

	input := "OOHHHH" // try different combinations

	for _, ch := range input {
		wg.Add(1)

		if ch == 'H' {
			go func() {
				defer wg.Done()
				h2o.Hydrogen(func() {
					fmt.Print("H")
				})
			}()
		} else {
			go func() {
				defer wg.Done()
				h2o.Oxygen(func() {
					fmt.Print("O")
				})
			}()
		}
	}

	wg.Wait()
}
