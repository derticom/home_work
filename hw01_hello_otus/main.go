package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	originalStr := "Hello, OTUS!"
	reversedStr := reverse.String(originalStr)

	fmt.Println(reversedStr)
}
