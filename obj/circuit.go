package obj

import (
	"errors"
	"log"
	"math"

	"github.com/ko6bxl/cm2go/block"
)

const NoOff = "noOff"

// Full Adder
// TODO: Make CLA Adder
func Add(bits int) (block.Collection, AddIO) {
	//Init variables
	var add AddIO
	var adder block.Collection
	//make each level of the adder for each bit
	for level := range bits {
		//add input 1
		in1 := adder.Append(block.NODE())
		in1.Offset.X = 0
		in1.Offset.Y = float32(level)
		//add input 2
		in2 := adder.Append(block.NODE())
		in2.Offset.X = 1
		in2.Offset.Y = float32(level)
		//add carry input
		in3 := adder.Append(block.OR())
		in3.Offset.Z = 0
		in3.Offset.X = 2
		in3.Offset.Y = float32(level)
		//add output
		out1 := adder.Append(block.NODE())
		out1.Offset.Z = 3
		out1.Offset.X = 1
		out1.Offset.Y = float32(level)
		//Logic
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
		//Carry out
		or1 := adder.Append(block.OR())
		or1.Offset.Z = 3
		or1.Offset.X = 0
		or1.Offset.Y = float32(level)
		//connect inputs to xor
		adder.Connect(in1, xor1)
		adder.Connect(in2, xor1)
		//connect carry to xor2 & and2
		adder.Connect(in3, xor2)
		adder.Connect(in3, and2)
		//connect inputs to and
		adder.Connect(in1, and1)
		adder.Connect(in2, and1)
		//connect internal logic
		adder.Connect(xor1, xor2)
		adder.Connect(xor1, and2)
		adder.Connect(and1, or1)
		adder.Connect(and2, or1)
		//connect output
		adder.Connect(xor2, out1)
		//append IO
		add.CIn = append(add.CIn, in3)
		add.AIn = append(add.AIn, in1)
		add.BIn = append(add.BIn, in2)
		add.COut = append(add.COut, or1)
		add.AOut = append(add.AOut, out1)
		//if above the first layer, chain carry's
		if level > 0 {
			adder.Connect(add.COut[level-1], add.CIn[level])
		}

	}

	return adder, add
}

// DEPRACATED:
// TODO: Remove this function lol
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

// Negate
func Negate(bits int) (block.Collection, NegateIO) {
	//init variables
	var negate block.Collection
	var nio NegateIO
	//generate fulladder
	add, aio := Add(bits)
	//Merge adder
	//TODO: Use new merge func
	negate, err := Merge(negate, add, "Z")
	if err != nil {
		log.Fatal(err)
	}
	//make flip flop for 2's complement
	flip := negate.Append(block.FLIPFLOP())
	flip.Offset.Y = 0
	flip.Offset.X = 1
	flip.Offset.Z = 0
	flip.State = true
	negate.Connect(flip, aio.BIn[0])
	//make inverters
	for i := range bits {
		nor := negate.Append(block.NOR())
		nor.Offset.Y = float32(i)
		nor.Offset.Z = 0

		negate.Connect(nor, aio.AIn[i])
		nio.AIn = append(nio.AIn, nor)
	}
	//adder output is the same as the negate
	nio.AOut = aio.AOut
	return negate, nio
}

// N bit decoder
func Decoder(bits int) (block.Collection, DecodeIO) {
	//init variables
	var decode block.Collection
	var decodeIO DecodeIO
	//keep track of non IO components
	var andPile []*block.Base
	var norPile []*block.Base
	var orPile []*block.Base
	//stores the Y level of the andpile
	up := 0
	//calculates the total ammount of and gates needed
	totalAnd := math.Pow(2, float64(bits))
	//casts into an int
	totalAndInt := int(totalAnd)
	//for input bits
	for i := range bits {
		//create or for passing whether a bit is 1
		or := decode.Append(block.OR())
		or.Offset.X = float32(i)
		or.Offset.Y = 0
		or.Offset.Z = -1
		//create nor for passing whether a bit is 0
		nor := decode.Append(block.NOR())
		nor.Offset.X = float32(i)
		nor.Offset.Y = 1
		nor.Offset.Z = -1
		//add node for input
		node := decode.Append(block.NODE())
		node.Offset.X = float32(i)
		node.Offset.Z = -2
		//connect input to or * nor
		decode.Connect(node, or)
		decode.Connect(node, nor)
		//append piles of non-io
		norPile = append(norPile, nor)
		orPile = append(orPile, or)
		//append io
		decodeIO.AIn = append(decodeIO.AIn, node)
	}
	//for 2^bits
	//uses outdated for loop to do some nasty shit
	for i := 0; i < totalAndInt; i++ {
		//create and for comparing 1 & 0 bits
		and := decode.Append(block.AND())
		and.Offset.X = float32(i)
		and.Offset.Y = float32(up)
		//append non-io piles
		andPile = append(andPile, and)
		//append IO
		decodeIO.AOut = append(decodeIO.AOut, and)
		//makes it so the and gates are translated up
		if i >= (totalAndInt/8)-1 {
			i = -1
			up++
		}
		//REALLY IMPORTANT: ends loop lol. without this, you get a fork bomb
		if up == 8 {
			break
		}
	}
	//do less horible iteration of 2^n bits
	//think of i as the number inputed into the decoder
	for i := range totalAndInt {
		//connect all nors when all inputs are 0
		if i == 0 {
			for _, nor := range norPile {
				decode.Connect(nor, andPile[0])
			}
		}
		//for n bits
		for x := range bits {
			//iterate through each bit in i
			shift := (i >> x)
			//checks if bit is 1 or 0 by checking if the shift made it odd or even
			if float32(shift)/2 != float32(shift/2) {
				//if 1
				decode.Connect(orPile[x], andPile[i])
			} else {
				//if 0
				decode.Connect(norPile[x], andPile[i])
			}
		}

	}
	return decode, decodeIO
}

// Subtractor
func Sub(bits int) (block.Collection, SubIO) {
	//init variables
	var sub block.Collection
	var subIO SubIO
	//create negate and adder
	negate, negIO := Negate(bits)
	add, addIO := Add(bits)
	//Merge the two
	//TODO: Use better merge function
	sub, err := Merge(negate, add, "Z")

	if err != nil {
		log.Fatal(err)
	}
	//connect negate out to adder in
	for i, aOut := range negIO.AOut {
		sub.Connect(aOut, addIO.BIn[i])
	}
	//make io
	subIO.AIn = addIO.AIn
	subIO.AOut = addIO.AOut
	subIO.COut = addIO.COut

	subIO.BIn = negIO.AIn

	return sub, subIO

}

func Mux(bits int) (block.Collection, MuxIO) {
	//init variables
	var mux block.Collection
	var muxIO MuxIO
	//create switch for mux
	ctrl := mux.Append(block.NODE())
	ctrl.Offset.X = -1
	//create inverse of that switch
	nor := mux.Append(block.NOR())
	nor.Offset.X = -1
	nor.Offset.Y = 1
	mux.Connect(ctrl, nor)
	//iterate through bits
	for i := range bits {
		//first byte switch
		andA := mux.Append(block.AND())
		andA.Offset.Y = float32(i)
		//second byte switch
		andB := mux.Append(block.AND())
		andB.Offset.Y = float32(i)
		andB.Offset.X = 1
		//output
		node := mux.Append(block.NODE())
		node.Offset.Y = float32(i)
		node.Offset.Z = 1
		//append io
		muxIO.AIn = append(muxIO.AIn, andA)
		muxIO.BIn = append(muxIO.BIn, andB)
		muxIO.AOut = append(muxIO.AOut, node)
		//connect logic
		mux.Connect(andA, node)
		mux.Connect(andB, node)
		mux.Connect(ctrl, andB)
		mux.Connect(nor, andA)
	}

	return mux, muxIO
}

// Creates bitwise and
func And(bits int) (block.Collection, AndIO) {
	//init variables
	var and block.Collection
	var andIO AndIO
	//iterate bits
	for i := range bits {
		//create and gates
		theAnd := and.Append(block.AND())
		theAnd.Offset.Y = float32(i)
		//append io
		andIO.AIn = append(andIO.AIn, theAnd)
		andIO.BIn = append(andIO.BIn, theAnd)
		andIO.AOut = append(andIO.AOut, theAnd)
	}

	return and, andIO
}

// creates a register
func Register(bits int) (block.Collection, RegIO) {
	//init variables
	var reg block.Collection
	var regIO RegIO
	//Write input
	ctrl := reg.Append(block.NODE())
	ctrl.Offset.Z = 1
	regIO.CIn = ctrl
	//iterate through bits
	for i := range bits {
		//flipflop to store the value
		flip := reg.Append(block.FLIPFLOP())
		flip.Offset.X = 1
		flip.Offset.Y = float32(i)
		//xor & and for control logic
		xor := reg.Append(block.XOR())
		xor.Offset.X = 1
		xor.Offset.Z = 1
		xor.Offset.Y = float32(i)

		and := reg.Append(block.AND())
		and.Offset.Y = float32(i)
		//connect logic
		reg.Connect(regIO.CIn, and)
		reg.Connect(xor, and)
		reg.Connect(flip, xor)
		reg.Connect(and, flip)
		//append io
		regIO.AIn = append(regIO.AIn, xor)
		regIO.AOut = append(regIO.AOut, flip)
	}
	return reg, regIO
}

// switch controls whether the input is outputed
func Switch(bits int) (block.Collection, SwitchIO) {
	//init io variable
	var switchIO SwitchIO
	//create bitwise and
	and, andIO := And(bits)
	//create switch control
	ctrl := and.Append(block.NODE())
	ctrl.Offset.X = 2
	//iterate bits
	for i := range bits {
		//connect control to and's
		and.Connect(ctrl, andIO.BIn[i])
	}
	//create io
	switchIO.AIn = andIO.AIn
	switchIO.CIn = ctrl
	switchIO.AOut = andIO.AOut
	return and, switchIO
}

// memory stores values in addresses and lets you read those values from addresses
// TODO: make separate func that lets creator choose ammount of addresses
func Memory(bits int) (block.Collection, MemIO) {
	//init variables
	var mem block.Collection
	var memIO MemIO
	//calculate number of registers
	var size int = int(math.Pow(2, float64(bits)))
	//controls writing to an address
	write := mem.Append(block.NODE())
	write.Offset.X = 2
	memIO.WIn = write
	//iterate bitdepth
	for i := range bits {
		//input for value
		in := mem.Append(block.NODE())
		in.Offset.Y = float32(i)
		memIO.AIn = append(memIO.AIn, in)
		//input for address
		adrs := mem.Append(block.NODE())
		adrs.Offset.X = 1
		adrs.Offset.Y = float32(i)
		memIO.BIn = append(memIO.BIn, adrs)
		//output of value in address
		out := mem.Append(block.NODE())
		out.Offset.X = 5
		out.Offset.Y = float32(i)
		memIO.AOut = append(memIO.AOut, out)
	}
	//decoder to choose address
	dec, decIO := Decoder(bits)
	//offset decoder
	for _, blk := range dec.Blocks {
		blk.Offset.X = 3
	}
	//merge decoder and memory IO created earlier
	mem.Merge(&dec)
	//iterate bitdepth
	for i := range bits {
		//connect address input to the decoder
		mem.Connect(memIO.BIn[i], decIO.AIn[i])
	}
	//init final register variables
	var finReg block.Collection
	var finRegIO finRegIO
	//for ammount of registers, = to 2^bits
	for i := range size {
		//create register and switch for output
		reg, regIO := Register(bits)
		switchh, switchIO := Switch(bits)
		//create write and
		wAnd, wAndIO := And(1)
		//iterate through write switch
		for _, blk := range wAnd.Blocks {
			blk.Offset.Z = 5
			blk.Offset.X = 2 + float32(i)
		}
		//iterate through register
		for _, blk := range reg.Blocks {
			blk.Offset.Z = blk.Offset.Z + 6
			blk.Offset.X = blk.Offset.X + 2 + float32(i*2)
		}
		//iterate through output switch
		for _, blk := range switchh.Blocks {
			blk.Offset.Z = blk.Offset.Z + 8
			blk.Offset.X = blk.Offset.X + 2 + float32(i)
		}
		//merge components
		reg.Merge(&switchh)
		reg.Merge(&wAnd)
		//connect write switch out to write
		reg.Connect(wAndIO.AOut[0], regIO.CIn)
		//for bitdepth
		for ii := range bits {
			//connect register output to output switch
			reg.Connect(regIO.AOut[ii], switchIO.AIn[ii])
		}
		//merge register to final register
		finReg.Merge(&reg)
		//append finalregister io
		finRegIO.Reg = append(finRegIO.Reg, regIO)
		finRegIO.RIn = append(finRegIO.RIn, switchIO.CIn)
		finRegIO.WIn = append(finRegIO.WIn, wAndIO.AIn[0])
		finRegIO.AOut = append(finRegIO.AOut, switchIO.AOut)
	}
	//merge final register
	mem.Merge(&finReg)
	//for 2^bits
	for i := range size {
		//connect decoder to read & write and
		mem.Connect(decIO.AOut[i], finRegIO.RIn[i])
		mem.Connect(decIO.AOut[i], finRegIO.WIn[i])
		//connect write to write and
		mem.Connect(memIO.WIn, finRegIO.WIn[i])
		//for the array inside of the array of outputs???
		for ii, x := range finRegIO.AOut[i] {
			//whatever, connect to memory output
			mem.Connect(x, memIO.AOut[ii])
		}
		//same thing but for inputs
		for ii, x := range finRegIO.Reg[i].AIn {
			mem.Connect(memIO.AIn[ii], x)
		}
	}

	return mem, memIO
}

// next few funcs are basically the same, so I am not going to comment them all
func Xor(bits int) (block.Collection, XorIO) {
	var xor block.Collection
	var xorIO XorIO

	for i := range bits {
		theXor := xor.Append(block.XOR())
		theXor.Offset.Y = float32(i)

		xorIO.AIn = append(xorIO.AIn, theXor)
		xorIO.BIn = append(xorIO.BIn, theXor)
		xorIO.AOut = append(xorIO.AOut, theXor)
	}
	return xor, xorIO
}

func Xnor(bits int) (block.Collection, XnorIO) {
	var xnor block.Collection
	var xnorIO XnorIO

	for i := range bits {
		theXnor := xnor.Append(block.XNOR())
		theXnor.Offset.Y = float32(i)

		xnorIO.AIn = append(xnorIO.AIn, theXnor)
		xnorIO.BIn = append(xnorIO.BIn, theXnor)
		xnorIO.AOut = append(xnorIO.AOut, theXnor)
	}
	return xnor, xnorIO
}

func Nor(bits int) (block.Collection, NorIO) {
	var nor block.Collection
	var norIO NorIO

	for i := range bits {
		thenor := nor.Append(block.NOR())
		thenor.Offset.Y = float32(i)

		norIO.AIn = append(norIO.AIn, thenor)
		norIO.BIn = append(norIO.BIn, thenor)
		norIO.AOut = append(norIO.AOut, thenor)
	}
	return nor, norIO
}

func Or(bits int) (block.Collection, OrIO) {
	var Or block.Collection
	var OrIO OrIO

	for i := range bits {
		theor := Or.Append(block.OR())
		theor.Offset.Y = float32(i)

		OrIO.AIn = append(OrIO.AIn, theor)
		OrIO.BIn = append(OrIO.BIn, theor)
		OrIO.AOut = append(OrIO.AOut, theor)
	}
	return Or, OrIO
}

func Nand(bits int) (block.Collection, NandIO) {
	var nand block.Collection
	var nandIO NandIO

	for i := range bits {
		thenand := nand.Append(block.NAND())
		thenand.Offset.Y = float32(i)

		nandIO.AIn = append(nandIO.AIn, thenand)
		nandIO.BIn = append(nandIO.BIn, thenand)
		nandIO.AOut = append(nandIO.AOut, thenand)
	}
	return nand, nandIO
}

func Led(bits int) (block.Collection, LedIO) {
	var led block.Collection
	var ledIO LedIO

	for i := range bits {
		leds := led.Append(block.LED(nil))
		leds.Offset.Y = float32(i)

		ledIO.AIn = append(ledIO.AIn, leds)
	}

	return led, ledIO
}

func Node(bits int) (block.Collection, NodeIO) {
	var node block.Collection
	var nodeIO NodeIO

	for i := range bits {
		nodes := node.Append(block.NODE())
		nodes.Offset.Y = float32(i)

		nodeIO.AIn = append(nodeIO.AIn, nodes)
		nodeIO.AOut = nodeIO.AIn
	}

	return node, nodeIO
}
