package protocol

type Header struct {
	Sender        *Client
	OPCode        OPCODE
	RequestLength int
}

func NewHeader(b []byte) (header *Header) {
	header = &Header{}

	header.OPCode = GetOPCode(b[0])

	return
}

func (h *Header) GetOPCode() OPCODE {
	return h.OPCode
}

func (h *Header) GetRequestLength() int {
	return h.RequestLength
}
