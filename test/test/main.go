package main

import (
	"fmt"
)

func main() {
	for i := range 30 {

		if float32(i)/2 != float32(i/2) {
			fmt.Println(i)
		}
	}
}
