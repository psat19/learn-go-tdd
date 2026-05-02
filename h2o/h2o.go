package main

import (
	"fmt"
	"sync"
)

type H2O struct {
	mu     sync.Mutex
	cond   *sync.Cond
	hCount int
	hStep  int
	oCount int
	oStep  int
}

func (b *H2O) AwaitHydrogen() {
	b.mu.Lock()

	for b.hStep == 0 {
		b.cond.Wait()
	}

	b.hCount++
	b.hStep--
	b.tryTrip()

	b.mu.Unlock()
}

func (b *H2O) AwaitOxygen() {
	b.mu.Lock()

	b.mu.Unlock()
}

func (b *H2O) tryTrip() {
	if b.hCount >= 2 && b.oCount >= 1 {
		b.hCount -= 2
		b.oCount -= 1
		b.hStep += 2
		b.oStep += 1
		b.cond.Broadcast()
	}
}

func NewH2O() *H2O {
	h2o := &H2O{}
	h2o.cond = sync.NewCond(&h2o.mu)

	return h2o
}

func (h2o *H2O) Hydrogen(releaseHydrogen func()) {
	h2o.AwaitHydrogen()
	releaseHydrogen()
}

func (h2o *H2O) Oxygen(releaseOxygen func()) {
	h2o.AwaitOxygen()
	releaseOxygen()
}

func main() {
	h2o := NewH2O()
	var wg sync.WaitGroup

	input := "HHHHOO" // try different combinations

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
