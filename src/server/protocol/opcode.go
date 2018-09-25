package protocol

type OPCODE int

const (
	PingFrame OPCODE = iota // 0
	AuthenticateFrame
	MessageFrame
	UnknownFrame
)

func GetOPCode(b byte) OPCODE {
	if OPCODE(b) >= UnknownFrame {
		return UnknownFrame
	}

	return OPCODE(b)
}
