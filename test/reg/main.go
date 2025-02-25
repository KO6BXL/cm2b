package main

import (
	"fmt"
	"log"

	"github.com/ko6bxl/cm2b/obj"
	"github.com/ko6bxl/cm2go/block"
	"github.com/ko6bxl/cm2go/build"
)

type MIO struct {
	AIn []*block.Base
	CIn *block.Base

	AOut []*block.Base
}

func main() {
	var test block.Collection
	var testIO MIO
	bits := 8
	reg, regIO := obj.Register(bits)

	//test, err := obj.Merge(test, reg, "Z")

	ctrl := test.Append(block.NODE())
	ctrl.Offset.Z = -5
	ctrl.Offset.X = 1

	testIO.CIn = ctrl

	for i := range bits {
		node := test.Append(block.FLIPFLOP())
		node.Offset.Z = -5
		node.Offset.Y = float32(i)

		led := test.Append(block.LED(nil))
		led.Offset.Z = -5
		led.Offset.X = 3
		led.Offset.Y = float32(i)

		testIO.AOut = append(testIO.AOut, led)
		testIO.AIn = append(testIO.AIn, node)

		test.Connect(testIO.AIn[i], regIO.AIn[i])
		test.Connect(regIO.AOut[i], testIO.AOut[i])
	}

	test.Connect(testIO.CIn, regIO.CIn)

	for i := range bits {
		log.Println("testIO")
		log.Println(testIO.AIn[i])
		log.Println("regIO")
		log.Println(regIO.AIn[i])
	}
	test, _ = obj.Merge(test, reg, "Z")
	fin, err := build.Compile([]block.Collection{test})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fin)
	log.Println(regIO)

}
