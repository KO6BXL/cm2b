package main

import (
	"fmt"
	"log"

	"github.com/ko6bxl/cm2b/cm2go/block"
	"github.com/ko6bxl/cm2b/cm2go/build"
	"github.com/ko6bxl/cm2b/obj"
)

func main() {

	test, _ := obj.Decoder(10)

	fin, err := build.Compile([]block.Collection{test})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fin)

}
