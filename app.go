package main

// This is a placeholder program, unrelated to the workshop itself. You
// are not expected to run or modify this file. Having a Go program in
// the project root avoids warning messages shown when running "go get"
// from this directory. See: https://github.com/golang/go/issues/37700

import (
	"fmt"
	"os"
	"log"
)

func main() {
	content, err := os.ReadFile("README.md")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(content))
}
