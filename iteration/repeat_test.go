package iteration

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	got := Repeat("a", 5)
	want := "aaaaa"

	if got != want {
		t.Errorf("got: %q \n want: %q", got, want)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a", 5)
	}
}

func ExampleRepeat() {
	got := Repeat("a", 5)
	fmt.Println(got)
	// Output: aaaaa
}
