package protocol

type OPCODE int

const (
	AuthenticateFrame OPCODE = iota
	MessageFrame
	UnknownFrame
)

func GetOPCode(b byte) OPCODE {
	switch b {
	case '1':
		return AuthenticateFrame
	case '2':
		return MessageFrame
	}

	return UnknownFrame
}
