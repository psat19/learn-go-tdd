package constants

import "fmt"

const spanish = "spanish"
const french = "french"

const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "

func Hello(name, language string) string {
	if name == "" {
		name = "World"
	}

	switch language {
	case spanish:
		return spanishHelloPrefix + name + "!!!"
	case french:
		return frenchHelloPrefix + name + "!!!"
	}

	return englishHelloPrefix + name + "!!!"
}

func main() {
	fmt.Println(Hello("Me", "english"))
}
