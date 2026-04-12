package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Me")
	want := "Hello, Me!!!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
