package obj

import (
	"github.com/ko6bxl/cm2go/block"
)

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

type DecodeIO struct {
	AIn  []*block.Base
	AOut []*block.Base
}

type SubIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
	COut []*block.Base
}

type MuxIO struct {
	AIn []*block.Base
	BIn []*block.Base
	CIn []*block.Base

	AOut []*block.Base
}

type AndIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

type RegIO struct {
	AIn []*block.Base
	CIn *block.Base

	AOut []*block.Base
}

type SwitchIO struct {
	AIn []*block.Base
	CIn *block.Base

	AOut []*block.Base
}

type MemIO struct {
	AIn []*block.Base
	BIn []*block.Base
	WIn *block.Base

	AOut []*block.Base
}

type finRegIO struct {
	Reg []RegIO
	RIn []*block.Base
	WIn []*block.Base

	AOut [][]*block.Base
}

type XorIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

type XnorIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

type NorIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

type OrIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

type NandIO struct {
	AIn []*block.Base
	BIn []*block.Base

	AOut []*block.Base
}

type LedIO struct {
	AIn []*block.Base
}
type NodeIO struct {
	AIn  []*block.Base
	AOut []*block.Base
}
