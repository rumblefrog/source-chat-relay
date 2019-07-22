package protocol

import "github.com/rumblefrog/source-chat-relay/server/packet"

type Deliverable interface {
	Marshal() []byte
	Content() string
	Author() string
}

type BaseMessage struct {
	Type MessageType

	// Internal relay purposes only
	SenderID string
}

func (b BaseMessage) Author() string {
	return b.SenderID
}

func ParseBaseMessage(r *packet.PacketReader) (m BaseMessage) {
	r.SetPos(0)

	m.Type = ParseMessageType(r.ReadUint8())

	return
}
