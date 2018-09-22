package protocol

import (
	"fmt"
	"strconv"
)

type Header struct {
	Sender        *Client
	OPCode        OPCODE
	RequestLength int
	PayloadLength int
}

func NewHeader(b []byte) (header *Header) {
	header = &Header{}

	header.OPCode = GetOPCode(b[0])

	header.PayloadLength = ParseLength(b[1:5])

	return
}

func ParseLength(b []byte) int {
	r, _ := strconv.ParseInt(fmt.Sprintf("%c%c%c%c", b[0], b[1], b[2], b[3]), 10, 0)

	return int(r)
}

func (h *Header) GetOPCode() OPCODE {
	return h.OPCode
}

func (h *Header) GetRequestLength() int {
	return h.RequestLength
}

func (h *Header) GetPayloadLength() int {
	return h.PayloadLength
}
