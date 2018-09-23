package protocol

import (
	"fmt"
	"strings"
)

const (
	HostnameLen   = 64
	ClientIDLen   = 64
	ClientNameLen = 32
)

type Message struct {
	Header     *Header
	Hostname   string
	ClientID   string
	ClientName string
	Content    string
}

func ParseMessage(b []byte, h *Header) *Message {
	offset := 2

	Message := &Message{}

	Message.Header = h

	for i := 0; i < HostnameLen; i++ {
		Message.Hostname += string(b[offset])
		offset++
	}

	for i := 0; i < ClientIDLen; i++ {
		Message.ClientID += string(b[offset])
		offset++
	}

	for i := 0; i < ClientNameLen; i++ {
		Message.ClientName += string(b[offset])
		offset++
	}

	for i := 0; i < h.GetRequestLength()-offset; i++ {
		Message.Content += string(b[offset])
		offset++
	}

	strings.TrimSpace(Message.Hostname)
	strings.TrimSpace(Message.ClientID)
	strings.TrimSpace(Message.ClientName)

	return Message
}

func (m *Message) ToString() (buffer string) {
	buffer += fmt.Sprintf("%s%-64s", buffer, m.Hostname)

	buffer += fmt.Sprintf("%s%-64s", buffer, m.ClientID)

	buffer += fmt.Sprintf("%s%-32s", buffer, m.ClientName)

	buffer += fmt.Sprintf("%s%s", buffer, m.Content)

	return
}

func (m *Message) GetHostname() string {
	return m.Hostname
}

func (m *Message) GetClientID() string {
	return m.ClientID
}

func (m *Message) GetClientName() string {
	return m.ClientName
}

func (m *Message) GetContent() string {
	return m.Content
}
