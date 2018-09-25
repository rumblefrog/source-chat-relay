package protocol

type OPCODE int

const (
	PingFrame         OPCODE = iota // 0
	AuthenticateFrame OPCODE = iota // 1
	MessageFrame      OPCODE = iota // 2
	UnknownFrame      OPCODE = iota
)

func GetOPCode(b byte) OPCODE {
	switch b {
	case '0':
		{
			return PingFrame
		}
	case '1':
		{
			return AuthenticateFrame
		}
	case '2':
		{
			return MessageFrame
		}
	}

	return UnknownFrame
}
