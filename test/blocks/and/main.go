package main

import (
	"fmt"
	"log"

	"github.com/ko6bxl/cm2b/obj"
	"github.com/ko6bxl/cm2go/block"
	"github.com/ko6bxl/cm2go/build"
)

func main() {

	test, _ := obj.And(64)

	fin, err := build.Compile([]block.Collection{test})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fin)

}
