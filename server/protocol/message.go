package protocol

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/server/helper"
	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
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
	message := ParseMessage(b, h)

	if message.Header.Sender.ID == "" {
		return
	}

	if message.Overwrite == nil {
		sender, err := repoEntity.GetEntity(message.Header.Sender.ID, repoEntity.Server)

		if err == nil && sender.DisplayName != message.Hostname {
			sender.DisplayName = message.Hostname
		}
	}

	log.WithFields(log.Fields{
		"Hostname":    message.Hostname,
		"Client ID":   message.ClientID,
		"Client Name": message.ClientName,
		"Content":     message.Content,
	}).Debug()

	NetManager.Router <- message
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

	buffer += fmt.Sprintf("%-64s", helper.StripSymbol(m.Hostname))

	buffer += fmt.Sprintf("%-64s", helper.StripSymbol(m.ClientID))

	buffer += fmt.Sprintf("%-32s", helper.StripSymbol(m.ClientName))

	buffer += fmt.Sprintf("%s", helper.StripSymbol(m.Content))

	return
}

func (m *Message) GetSendChannels() []int {
	if m.Overwrite != nil {
		return m.Overwrite.SendChannels
	}

	entity, err := repoEntity.GetEntity(m.Header.Sender.ID, repoEntity.Server)

	if err != nil {
		return []int{}
	}

	return entity.SendChannels
}

func (m *Message) GetReceiveChannels() []int {
	if m.Overwrite != nil {
		return m.Overwrite.ReceiveChannels
	}

	entity, err := repoEntity.GetEntity(m.Header.Sender.ID, repoEntity.Server)

	if err != nil {
		return []int{}
	}

	return entity.ReceiveChannels
}

func (m *Message) GetClientURL() string {
	return fmt.Sprintf("https://steamcommunity.com/profiles/%s", m.ClientID)
}

func (m *Message) GetClientColor() int {
	c := []byte(m.ClientID)

	i, _ := strconv.ParseInt(string(c[len(c)-6:]), 16, 64)

	return int(i)
}
