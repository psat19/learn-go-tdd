package main

import "github.com/psat19/learn-go-tdd/struct-interfaces/shapes"

func Perimeter(rect shapes.Rectangle) float64 {
	return 2 * (rect.Height + rect.Width)
}

func Area(rect shapes.Rectangle) float64 {
	return rect.Height * rect.Width
}
