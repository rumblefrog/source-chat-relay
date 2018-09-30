package protocol

type OPCODE int

const (
	PingFrame OPCODE = iota // 0
	AuthenticateFrame
	MessageFrame
	UnknownFrame
)

func GetOPCode(b byte) OPCODE {
	switch b {
	case '0':
		return PingFrame
	case '1':
		return AuthenticateFrame
	case '2':
		return MessageFrame
	}

	return UnknownFrame
}
