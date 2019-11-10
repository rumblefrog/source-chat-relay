package protocol

import (
	"strings"

	"github.com/rumblefrog/source-chat-relay/server/packet"
)

type KV struct {
	Name  string
	Value string
}

type KVMessage struct {
	BaseMessage

	Name string

	Count uint8

	KeyValues []KV
}

func ParseKVMessage(base BaseMessage, r *packet.PacketReader) (*KVMessage, error) {
	m := &KVMessage{}

	m.BaseMessage = base

	var ok bool

	m.Name, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	if err := r.CanRead(1); err != nil {
		return nil, ErrOutOfBound
	}

	m.Count = r.ReadUint8()

	for i := 0; i < int(m.Count); i++ {
		var kv KV

		kv.Name, ok = r.TryReadString()

		if !ok {
			return nil, ErrCannotReadString
		}

		kv.Value, ok = r.TryReadString()

		if !ok {
			return nil, ErrCannotReadString
		}

		m.KeyValues = append(m.KeyValues, kv)
	}

	return m, nil
}

func (m *KVMessage) Type() MessageType {
	return MessageKV
}

func (m *KVMessage) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(MessageKV))
	builder.WriteCString(m.BaseMessage.EntityName)

	builder.WriteCString(m.Name)
	builder.WriteByte(m.Count)

	for _, kv := range m.KeyValues {
		builder.WriteCString(kv.Name)
		builder.WriteCString(kv.Value)
	}

	return builder.Bytes()
}

func (m *KVMessage) Content() string {
	var builder strings.Builder

	for _, kv := range m.KeyValues {
		builder.WriteString(kv.Name + ":" + kv.Value + " ")
	}

	return builder.String()
}

func (m *KVMessage) Plain() string {
	return m.Content()
}

// TODO: Embed func
