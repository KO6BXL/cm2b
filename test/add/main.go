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
	//negate, _ := obj.Negate(4)
	// fin, err := build.FastCompile([]block.Collection{negate})
	var test block.Collection
	var inyes in
	add, addIO := obj.Add(8)

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

	test, err := obj.Merge(test, add, obj.NoOff)

	if err != nil {
		log.Fatal(err)
	}

	for i := range 8 {
		test.Connect(inyes.AIn[i], addIO.BIn[i])

		test.Connect(inyes.BIn[i], addIO.AIn[i])
	}

	//var test block.Collection
	//flip := test.Append(block.FLIPFLOP())
	//flip.State = true
	//flip.Offset.Y = 0
	//flip.Offset.X = 0
	//flip.Offset.Z = 2

	fin, err := build.Compile([]block.Collection{test})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fin)

}
