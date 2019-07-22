package protocol

import (
	"encoding/binary"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/packet"
)

type IdentificationType uint8

const (
	IdentificationInvalid IdentificationType = iota
	IdentificationSteam
	IdentificationDiscord
	IdentificationTypeCount
)

type ChatMessage struct {
	BaseMessage

	IDType IdentificationType

	ID string

	Username string

	Message string
}

func ParseChatMessage(base BaseMessage, r *packet.PacketReader) (m *ChatMessage) {
	m.BaseMessage = base

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
	builder.WriteByte(byte(m.IDType))
	builder.WriteCString(m.ID)
	builder.WriteCString(m.Username)
	builder.WriteCString(m.Message)

	return builder.Bytes()
}

func (m *ChatMessage) Plain() string {
	return m.Username + ": " + m.Message
}

func (m *ChatMessage) Embed() *discordgo.MessageEmbed {
	idColorBytes := []byte(m.ID)

	// Convert to an int with length of 6
	color := int(binary.BigEndian.Uint32(idColorBytes[len(idColorBytes)-6:])) / 10000

	return &discordgo.MessageEmbed{
		Color:       color,
		Description: m.Message,
		Timestamp:   time.Now().Format(time.RFC3339),
		Author: &discordgo.MessageEmbedAuthor{
			Name: m.Username,
			URL:  m.IDType.FormatURL(m.ID),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: m.BaseMessage.Hostname,
		},
	}
}

func ParseIdentificationType(t uint8) IdentificationType {
	if t >= uint8(IdentificationTypeCount) {
		return IdentificationInvalid
	}

	return IdentificationType(t)
}

func (i IdentificationType) FormatURL(id string) string {
	switch i {
	case IdentificationSteam:
		return "https://steamcommunity.com/profiles/" + id
	default:
		return ""
	}
}
