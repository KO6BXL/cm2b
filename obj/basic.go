package obj

import (
	"errors"
	"log"
	"math"

	"github.com/ko6bxl/cm2b/cm2go/block"
)

const NoOff = "noOff"

func Add(bits int) (block.Collection, AddIO) {
	var add AddIO
	var adder block.Collection

	for level := range bits {
		in1 := adder.Append(block.NODE())
		in1.Offset.X = 0
		in1.Offset.Y = float32(level)

		in2 := adder.Append(block.NODE())
		in2.Offset.X = 1
		in2.Offset.Y = float32(level)

		in3 := adder.Append(block.OR())
		in3.Offset.Z = 0
		in3.Offset.X = 2
		in3.Offset.Y = float32(level)

		out1 := adder.Append(block.NODE())
		out1.Offset.Z = 3
		out1.Offset.X = 1
		out1.Offset.Y = float32(level)

		xor1 := adder.Append(block.XOR())
		xor1.Offset.Z = 1
		xor1.Offset.X = 0
		xor1.Offset.Y = float32(level)

		xor2 := adder.Append(block.XOR())
		xor2.Offset.Z = 1
		xor2.Offset.X = 1
		xor2.Offset.Y = float32(level)

		and1 := adder.Append(block.AND())
		and1.Offset.Z = 2
		and1.Offset.X = 0
		and1.Offset.Y = float32(level)

		and2 := adder.Append(block.AND())
		and2.Offset.Z = 2
		and2.Offset.X = 1
		and2.Offset.Y = float32(level)

		or1 := adder.Append(block.OR())
		or1.Offset.Z = 3
		or1.Offset.X = 0
		or1.Offset.Y = float32(level)

		adder.Connect(in1, xor1)
		adder.Connect(in2, xor1)
		adder.Connect(in3, xor2)
		adder.Connect(in3, and2)

		adder.Connect(in1, and1)
		adder.Connect(in2, and1)

		adder.Connect(xor1, xor2)
		adder.Connect(xor1, and2)
		adder.Connect(and1, or1)
		adder.Connect(and2, or1)

		adder.Connect(xor2, out1)

		add.CIn = append(add.CIn, in3)
		add.AIn = append(add.AIn, in1)
		add.BIn = append(add.BIn, in2)

		add.COut = append(add.COut, or1)
		add.AOut = append(add.AOut, out1)
		if level > 0 {
			adder.Connect(add.COut[level-1], add.CIn[level])
		}

	}
	return adder, add
}

func Merge(col1, col2 block.Collection, direction string) (block.Collection, error) {
	var new block.Collection
	var mostLength int
	switch direction {
	case "Z":
		for _, block := range col1.Blocks {
			if mostLength < int(block.Offset.Z) {
				mostLength = int(block.Offset.Z)
			}
		}

		for _, block := range col2.Blocks {
			block.Offset.Z += float32(mostLength) + 1
		}
	case "Y":
		for _, block := range col1.Blocks {
			if mostLength < int(block.Offset.Y) {
				mostLength = int(block.Offset.Y)
			}
		}
		for _, block := range col2.Blocks {
			block.Offset.Y += float32(mostLength) + 1
		}
	case "X":
		for _, block := range col1.Blocks {
			if mostLength < int(block.Offset.X) {
				mostLength = int(block.Offset.X)
			}
		}
		for _, block := range col2.Blocks {
			block.Offset.X += float32(mostLength) + 1
		}
	case "noOff":
		log.Println("Declared 'noOff' in obj.Merge()")
	default:
		return new, errors.New("Direction not known!")
	}

	new.Blocks = append(col1.Blocks, col2.Blocks...)
	new.Connections = append(col1.Connections, col2.Connections...)
	return new, nil
}

type AddIO struct {
	CIn []*block.Base
	AIn []*block.Base
	BIn []*block.Base

	COut []*block.Base
	AOut []*block.Base
}

type NegateIO struct {
	AIn  []*block.Base
	AOut []*block.Base
}

func Negate(bits int) (block.Collection, NegateIO) {
	var negate block.Collection
	var nio NegateIO

	add, aio := Add(bits)

	negate, err := Merge(negate, add, "Z")

	if err != nil {
		log.Fatal(err)
	}

	flip := negate.Append(block.FLIPFLOP())
	flip.Offset.Y = 0
	flip.Offset.X = 1
	flip.Offset.Z = 0
	flip.State = true

	negate.Connect(flip, aio.BIn[0])

	for i := range bits {
		nor := negate.Append(block.NOR())
		nor.Offset.Y = float32(i)
		nor.Offset.Z = 0

		negate.Connect(nor, aio.AIn[i])
		nio.AIn = append(nio.AIn, nor)
	}

	nio.AOut = aio.AOut

	return negate, nio
}

type DecodeIO struct {
	AIn  []*block.Base
	AOut []*block.Base
}

func Decoder(bits int) (block.Collection, DecodeIO) {
	var decode block.Collection
	var decodeIO DecodeIO
	var andPile []*block.Base
	var norPile []*block.Base
	var orPile []*block.Base
	up := 0
	totalAnd := math.Pow(2, float64(bits))
	totalAndInt := int(totalAnd)

	for i := range bits {
		or := decode.Append(block.OR())
		or.Offset.X = float32(i)
		or.Offset.Y = 0
		or.Offset.Z = -1

		nor := decode.Append(block.NOR())
		nor.Offset.X = float32(i)
		nor.Offset.Y = 1
		nor.Offset.Z = -1

		node := decode.Append(block.NODE())
		node.Offset.X = float32(i)
		node.Offset.Z = -2

		decode.Connect(node, or)
		decode.Connect(node, nor)

		norPile = append(norPile, nor)
		orPile = append(orPile, or)
		decodeIO.AIn = append(decodeIO.AIn, node)
	}

	for i := 0; i < totalAndInt; i++ {
		and := decode.Append(block.AND())
		and.Offset.X = float32(i)
		and.Offset.Y = float32(up)

		andPile = append(andPile, and)
		decodeIO.AOut = append(decodeIO.AOut, and)

		if i >= (totalAndInt/8)-1 {
			i = -1
			up++
		}
		if up == 8 {
			break
		}
	}

	for i := range totalAndInt {
		if i == 0 {
			for _, nor := range norPile {
				decode.Connect(nor, andPile[0])
			}
		}
		for x := range bits {
			shift := (i >> x)
			if float32(shift)/2 != float32(shift/2) {
				decode.Connect(orPile[x], andPile[i])
			} else {
				decode.Connect(norPile[x], andPile[i])
			}
		}

	}

	return decode, decodeIO
}

type SubIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
	COut []*block.Base
}

func Sub(bits int) (block.Collection, SubIO) {
	var sub block.Collection
	var subIO SubIO
	negate, negIO := Negate(bits)
	add, addIO := Add(bits)

	sub, err := Merge(negate, add, "Z")

	if err != nil {
		log.Fatal(err)
	}

	for i, aOut := range negIO.AOut {
		sub.Connect(aOut, addIO.BIn[i])
	}

	subIO.AIn = addIO.AIn
	subIO.AOut = addIO.AOut
	subIO.COut = addIO.COut

	subIO.BIn = negIO.AIn

	return sub, subIO

}

type MuxIO struct {
	AIn []*block.Base
	BIn []*block.Base
	CIn []*block.Base

	AOut []*block.Base
}

func Mux(bits int) (block.Collection, MuxIO) {
	var mux block.Collection
	var muxIO MuxIO

	ctrl := mux.Append(block.NODE())
	ctrl.Offset.X = -1
	nor := mux.Append(block.NOR())
	nor.Offset.X = -1
	nor.Offset.Y = 1
	mux.Connect(ctrl, nor)
	for i := range bits {
		andA := mux.Append(block.AND())
		andA.Offset.Y = float32(i)
		andB := mux.Append(block.AND())
		andB.Offset.Y = float32(i)
		andB.Offset.X = 1
		node := mux.Append(block.NODE())
		node.Offset.Y = float32(i)
		node.Offset.Z = 1

		muxIO.AIn = append(muxIO.AIn, andA)
		muxIO.BIn = append(muxIO.BIn, andB)
		muxIO.AOut = append(muxIO.AOut, node)

		mux.Connect(andA, node)
		mux.Connect(andB, node)
		mux.Connect(ctrl, andB)
		mux.Connect(nor, andA)
	}

	return mux, muxIO
}

type AndIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

func And(bits int) (block.Collection, AndIO) {
	var and block.Collection
	var andIO AndIO

	for i := range bits {
		theAnd := and.Append(block.AND())
		theAnd.Offset.Y = float32(i)

		andIO.AIn = append(andIO.AIn, theAnd)
		andIO.BIn = append(andIO.BIn, theAnd)
		andIO.AOut = append(andIO.AOut, theAnd)
	}

	return and, andIO
}
