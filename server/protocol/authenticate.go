package protocol

import "github.com/rumblefrog/source-chat-relay/server/packet"

type AuthenticateMessage struct {
	BaseMessage

	Hostname string

	Token string
}

type AuthenticateMessageResponse struct {
	BaseMessage

	Response AuthenticateResponse
}

func ParseAuthenticateMessage(base BaseMessage, r *packet.PacketReader) (m *AuthenticateMessage) {
	m.BaseMessage = base

	m.Hostname = r.ReadString()

	m.Token = r.ReadString()

	return
}

// No marshal for authenticate message as we are the server and would never use it
// No parse authenticate message response because we should never receive it

func (m *AuthenticateMessageResponse) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(m.BaseMessage.Type))
	builder.WriteByte(byte(m.Response))

	return builder.Bytes()
}
