package protocol

import "github.com/rumblefrog/source-chat-relay/server/packet"

type ChatMessage struct {
	BaseMessage

	EntityName string

	IDType IdentificationType

	ID string

	Username string

	Message string
}

func ParseChatMessage(base BaseMessage, r *packet.PacketReader) (m *ChatMessage) {
	m.BaseMessage = base

	m.EntityName = r.ReadString()

	m.IDType = ParseIdentificationType(r.ReadUint8())

	m.ID = r.ReadString()

	m.Username = r.ReadString()

	m.Message = r.ReadString()

	return
}

func (m *ChatMessage) Content() string {
	return m.Message
}

func (m *ChatMessage) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(m.BaseMessage.Type))
	builder.WriteCString(m.EntityName)
	builder.WriteByte(byte(m.IDType))
	builder.WriteCString(m.ID)
	builder.WriteCString(m.Username)
	builder.WriteCString(m.Message)

	return builder.Bytes()
}
