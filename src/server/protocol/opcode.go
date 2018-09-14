package protocol

type OPCODE int

const (
	MessageFrame OPCODE = iota
	PingFrame    OPCODE = iota
	UnknownFrame OPCODE = iota
)

func GetOPCode(b byte) OPCODE {
	switch b {
	case '0':
		{
			return PingFrame
		}
	case '1':
		{
			return MessageFrame
		}
	}

	return UnknownFrame
}
