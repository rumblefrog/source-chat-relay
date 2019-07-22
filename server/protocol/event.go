package protocol

import "github.com/rumblefrog/source-chat-relay/server/packet"

type EventMessage struct {
	BaseMessage

	Event string

	Data string
}

func ParseEventMessage(base BaseMessage, r *packet.PacketReader) (m *EventMessage) {
	m.BaseMessage = base

	m.Event = r.ReadString()

	m.Data = r.ReadString()

	return
}

func (m *EventMessage) Content() string {
	return m.Data
}

func (m *EventMessage) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(m.BaseMessage.Type))
	builder.WriteString(m.Event)
	builder.WriteString(m.Data)

	return builder.Bytes()
}
