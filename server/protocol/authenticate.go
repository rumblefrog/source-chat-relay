package protocol

import (
	"github.com/rumblefrog/source-chat-relay/server/packet"
)

type AuthenticateMessage struct {
	BaseMessage

	Token string
}

type AuthenticateMessageResponse struct {
	BaseMessage

	Response AuthenticateResponse
}

func ParseAuthenticateMessage(base BaseMessage, r *packet.PacketReader) (*AuthenticateMessage, error) {
	m := &AuthenticateMessage{}

	m.BaseMessage = base

	var ok bool

	m.Token, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	return m, nil
}

// No marshal for authenticate message as we are the server and would never use it
// No parse authenticate message response because we should never receive it

func (m *AuthenticateMessageResponse) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(MessageAuthenticateResponse))
	builder.WriteCString("RELAY")

	builder.WriteByte(byte(m.Response))

	return builder.Bytes()
}
