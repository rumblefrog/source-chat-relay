package protocol

import (
	"strconv"
)

type Header struct {
	OPCode OPCODE
	Length uint16
}

func NewHeader(b []byte) (header *Header) {
	header = &Header{}

	header.OPCode = GetOPCode(b[0])

	header.Length = ParseLength(b[1:5])

	return
}

func ParseLength(b []byte) uint16 {
	r, _ := strconv.ParseUint(string(b[0]+b[1]+b[2]+b[3]), 10, 16)

	return uint16(r)
}
