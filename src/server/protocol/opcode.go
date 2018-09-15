package protocol

type OPCODE int

const (
	PingFrame    OPCODE = 1
	MessageFrame OPCODE = 5
	UnknownFrame OPCODE = 0
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
