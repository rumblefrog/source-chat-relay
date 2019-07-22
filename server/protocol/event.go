package protocol

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/packet"
)

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

func (m *EventMessage) Plain() string {
	return m.Event + ": " + m.Data
}

func (m *EventMessage) Embed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     16777215,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: m.BaseMessage.Hostname,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  m.Event,
				Value: m.Data,
			},
		},
	}
}
