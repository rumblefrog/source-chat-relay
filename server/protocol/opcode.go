package protocol

type OPCODE int

const (
	AuthenticateFrame OPCODE = iota // 0
	MessageFrame
	UnknownFrame
)

func GetOPCode(b byte) OPCODE {
	switch b {
	case '0':
		return AuthenticateFrame
	case '1':
		return MessageFrame
	}

	return UnknownFrame
}
