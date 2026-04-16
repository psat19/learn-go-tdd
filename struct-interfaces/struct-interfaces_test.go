package main

import (
	"testing"

	"github.com/psat19/learn-go-tdd/struct-interfaces/shapes"
)

func TestPerimeter(t *testing.T) {
	rect := shapes.Rectangle{Width: 10.0, Height: 10.0}
	got := Perimeter(rect)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	areaTests := []struct {
		name  string
		shape shapes.Shape
		want  float64
	}{
		{"Rectangle", shapes.Rectangle{Width: 12, Height: 6}, 72.0},
		{"Circle", shapes.Circle{Radius: 10}, 314.1592653589793},
		{"Triangle", shapes.Triangle{Base: 12, Height: 6}, 36.0},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.want {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.want)
			}
		})
	}

}
