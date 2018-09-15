package protocol

import (
	"log"
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
	var (
		hostname   string
		clientid   string
		clientname string
		content    string
		offset     = 5
	)

	for i := 0; i < HostnameLen; i++ {
		hostname += string(b[offset])
		offset++
	}

	log.Println(string(hostname))

	for i := 0; i < ClientIDLen; i++ {
		clientid += string(b[offset])
		offset++
	}

	for i := 0; i < ClientNameLen; i++ {
		clientname += string(b[offset])
		offset++
	}

	for i := 0; i < h.GetPayloadLength()-offset; i++ {
		content += string(b[offset])
		offset++
	}

	return &Message{
		Header:     h,
		Hostname:   strings.TrimSpace(hostname),
		ClientID:   strings.TrimSpace(clientid),
		ClientName: strings.TrimSpace(clientname),
		Content:    string(content),
	}
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
