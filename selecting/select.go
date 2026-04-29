package main

import (
	"errors"
	"net/http"
	"time"
)

func Racer(a, b string) (winner string, err error) {
	select {
	case <-ping(a):
		return a, nil

	case <-ping(b):
		return b, nil

	case <-time.After(10 * time.Second):
		return "", errors.New("Timeout")
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
