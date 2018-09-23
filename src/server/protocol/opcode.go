package protocol

type OPCODE int

const (
	PingFrame              OPCODE = iota // 0
	AuthenticateFrame      OPCODE = iota // 1
	AuthenticateReplyFrame OPCODE = iota // 2
	MessageFrame           OPCODE = iota // 3
	UnknownFrame           OPCODE = iota
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
			return AuthenticateReplyFrame
		}
	case '3':
		{
			return MessageFrame
		}
	}

	return UnknownFrame
}
