package main

import (
	"fmt"
	"log"

	"github.com/ko6bxl/cm2b/obj"
	"github.com/ko6bxl/cm2go/block"
	"github.com/ko6bxl/cm2go/build"
)

type in struct {
	AIn []*block.Base
	BIn []*block.Base
}

func main() {
	var test block.Collection
	var inyes in
	sub, subIO := obj.Sub(8)

	for i := range 8 {
		node := test.Append(block.NODE())
		node.Offset.Z = -5
		node.Offset.Y = float32(i)
		inyes.AIn = append(inyes.AIn, node)

	}
	for i := range 8 {
		node := test.Append(block.NODE())
		node.Offset.Z = -5
		node.Offset.X = 1
		node.Offset.Y = float32(i)

		inyes.BIn = append(inyes.BIn, node)
	}

	test, err := obj.Merge(test, sub, obj.NoOff)

	if err != nil {
		log.Fatal(err)
	}

	for i := range 8 {
		test.Connect(inyes.AIn[i], subIO.BIn[i])

		test.Connect(inyes.BIn[i], subIO.AIn[i])
	}

	fin, err := build.Compile([]block.Collection{test})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fin)

}
