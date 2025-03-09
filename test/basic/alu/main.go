package main

import (
	"fmt"
	"log"

	"github.com/ko6bxl/cm2b/obj"
	"github.com/ko6bxl/cm2go/block"
	"github.com/ko6bxl/cm2go/build"
)

type MIn struct {
	AIn []*block.Base
	BIn []*block.Base
	CIn []*block.Base

	AOut []*block.Base
}

func main() {
	var alu block.Collection
	var mIn MIn
	const bitDepth = 512

	add, addIO := obj.Add(bitDepth)

	sub, subIO := obj.Sub(bitDepth)

	dec, decIO := obj.Decoder(2)

	AAnd, AAndIO := obj.And(bitDepth)

	BAnd, BAndIO := obj.And(bitDepth)

	for i := range bitDepth {
		inA := alu.Append(block.FLIPFLOP())
		inA.Offset.Y = float32(i)
		inA.Offset.Z = -5

		inB := alu.Append(block.FLIPFLOP())
		inB.Offset.X = 1
		inB.Offset.Y = float32(i)
		inB.Offset.Z = -5

		outA := alu.Append(block.NODE())
		outA.Offset.Y = float32(i)
		outA.Offset.X = 5
		outA.Offset.Z = -5

		mIn.AIn = append(mIn.AIn, inA)
		mIn.BIn = append(mIn.BIn, inB)
		mIn.AOut = append(mIn.AOut, outA)
	}

	alu, _ = obj.Merge(alu, dec, "X")
	alu, _ = obj.Merge(alu, AAnd, "X")
	alu, _ = obj.Merge(alu, BAnd, "X")

	for i := range 2 {
		inC := alu.Append(block.FLIPFLOP())
		inC.Offset.Y = float32(i)
		inC.Offset.X = 2
		inC.Offset.Z = -5

		mIn.CIn = append(mIn.CIn, inC)

		alu.Connect(mIn.CIn[i], decIO.AIn[i])
	}

	alu, _ = obj.Merge(alu, add, "Z")
	alu, _ = obj.Merge(alu, sub, "Z")

	for i := range bitDepth {
		alu.Connect(decIO.AOut[1], AAndIO.AIn[i])
		alu.Connect(decIO.AOut[2], BAndIO.AIn[i])

		alu.Connect(addIO.AOut[i], AAndIO.BIn[i])
		alu.Connect(subIO.AOut[i], BAndIO.BIn[i])

		alu.Connect(mIn.AIn[i], addIO.AIn[i])
		alu.Connect(mIn.AIn[i], subIO.BIn[i])
		alu.Connect(mIn.BIn[i], addIO.BIn[i])
		alu.Connect(mIn.BIn[i], subIO.AIn[i])

		alu.Connect(AAndIO.AOut[i], mIn.AOut[i])
		alu.Connect(BAndIO.AOut[i], mIn.AOut[i])
	}

	fin, err := build.Compile([]block.Collection{alu})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fin)

}
