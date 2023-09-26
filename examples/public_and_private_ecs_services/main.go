package main

import (
	"fmt"
	"salami/common/driver"
)

func main() {
	errors := driver.Run()

	for _, err := range errors {
		fmt.Println(err.Error())
	}
}
