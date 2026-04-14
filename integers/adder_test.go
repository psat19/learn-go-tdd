package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	t.Run(
		"should add two numbers together",
		func(t *testing.T) {
			got := Add(2, 3)
			want := 0
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)

}
