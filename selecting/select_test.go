package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeTestServerWithDelay(20 * time.Millisecond)
		fastServer := makeTestServerWithDelay(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, _ := Racer(slowURL, fastURL)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		serverA := makeTestServerWithDelay(11 * time.Second)
		serverB := makeTestServerWithDelay(12 * time.Second)

		defer serverA.Close()
		defer serverB.Close()

		_, err := Racer(serverA.URL, serverB.URL)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}

func makeTestServerWithDelay(d time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if d != 0 {
		// 	time.Sleep(d)
		// }
		// w.Header().Write(os.Stdout) - Check what this does.

		if d != 0 {
			time.Sleep(d)
		}
		w.WriteHeader(200)
	}))

}
