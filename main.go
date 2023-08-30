package main

import (
	"fmt"
	"salami/compiler"
)

func main() {
	if err := compiler.Compile(); err != nil {
		fmt.Printf("Compilation failed: %v\n", err)
	}
}
