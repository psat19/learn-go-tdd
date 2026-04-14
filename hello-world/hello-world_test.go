package constants

import "testing"

func TestHello(t *testing.T) {
	t.Run(
		"When the hello function has variable, it has to spit out that variable",
		func(t *testing.T) {
			got := Hello("Me", "")
			want := "Hello, Me!!!"

			assertCorrectMessage(t, got, want)
		})

	t.Run(
		"Empty string should default to world",
		func(t *testing.T) {
			got := Hello("", "")
			want := "Hello, World!!!"

			assertCorrectMessage(t, got, want)
		})

	t.Run(
		"Empty string should default to world and spanish should be hola",
		func(t *testing.T) {
			got := Hello("", "spanish")
			want := "Hola, World!!!"

			assertCorrectMessage(t, got, want)
		})

	t.Run(
		"Me and spanish should be hola, Me",
		func(t *testing.T) {
			got := Hello("Me", "spanish")
			want := "Hola, Me!!!"

			assertCorrectMessage(t, got, want)
		})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	// t.Helper() marks the calling function as a test helper function. When printing file and line information, that function will be skipped. This allows us to see the line of our test code where the failure occurred, rather than inside the assertCorrectMessage function.
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
