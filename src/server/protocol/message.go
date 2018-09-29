package protocol

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	HostnameLen   = 64
	ClientIDLen   = 64
	ClientNameLen = 32
)

type Message struct {
	Header     *Header
	Overwrite  *OverwriteData
	Hostname   string
	ClientID   string
	ClientName string
	Content    string
}

type OverwriteData struct {
	ReceiveChannels []int
	SendChannels    []int
}

func (m *ClientManager) HandleMessage(b []byte, h *Header) {
	Message := ParseMessage(b, h)

	if Message.Header.Sender.Entity == nil {
		return
	}

	log.WithFields(log.Fields{
		"Hostname":      Message.Hostname,
		"Client ID":     Message.ClientID,
		"Client Name":   Message.ClientName,
		"Content":       Message.Content,
		"Send Channels": Message.Header.Sender.Entity.SendChannels,
	}).Debug()

	NetManager.Router <- Message
}

func ParseMessage(b []byte, h *Header) *Message {
	offset := 1

	Message := &Message{}

	Message.Header = h

	Message.Overwrite = nil

	Message.Hostname = string(b[offset : offset+HostnameLen])

	offset += HostnameLen

	Message.ClientID = string(b[offset : offset+ClientIDLen])

	offset += ClientIDLen

	Message.ClientName = string(b[offset : offset+ClientNameLen])

	offset += ClientNameLen

	Message.Content = string(b[offset:])

	Message.Hostname = strings.TrimSpace(Message.Hostname)
	Message.ClientID = strings.TrimSpace(Message.ClientID)
	Message.ClientName = strings.TrimSpace(Message.ClientName)

	return Message
}

func (m *Message) ToString() (buffer string) {
	buffer += "1"

	buffer += fmt.Sprintf("%-64s", m.Hostname)

	buffer += fmt.Sprintf("%-64s", m.ClientID)

	buffer += fmt.Sprintf("%-32s", m.ClientName)

	buffer += fmt.Sprintf("%s", m.Content)

	return
}

func (m *Message) GetSendChannels() []int {
	if m.Overwrite != nil {
		return m.Overwrite.SendChannels
	}

	return m.Header.Sender.Entity.SendChannels
}

func (m *Message) GetReceiveChannels() []int {
	if m.Overwrite != nil {
		return m.Overwrite.ReceiveChannels
	}

	return m.Header.Sender.Entity.ReceiveChannels
}

func (m *Message) GetClientURL() string {
	return fmt.Sprintf("https://steamcommunity.com/profiles/%s", m.ClientID)
}

func (m *Message) GetClientColor() int {
	c := []byte(m.ClientID)

	i, _ := strconv.ParseInt(string(c[len(c)-6:]), 16, 64)

	return int(i)
}
