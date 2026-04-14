package iteration

import "strings"

func Repeat(input string, times int) string {
	var builder strings.Builder
	builder.Grow(times)

	for i := 0; i < times; i++ {
		builder.WriteString(input)
	}

	return builder.String()
}
