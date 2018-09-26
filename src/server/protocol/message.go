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
	Hostname   string
	ClientID   string
	ClientName string
	Content    string
}

func (m *ClientManager) HandleMessage(b []byte, h *Header) {
	Message := ParseMessage(b, h)

	log.Println(Message.ClientName)

	// TODO: Send to router & bot chan

	// log.Printf("Hostname: %s \n", Message.Hostname)
	// log.Printf("ID: %s \n", Message.ClientID)
	// log.Printf("Name: %s \n", Message.ClientName)
	// log.Printf("Content: %s", Message.Content)
}

func ParseMessage(b []byte, h *Header) *Message {
	offset := 2

	Message := &Message{}

	Message.Header = h

	Message.Hostname = string(b[offset : offset+HostnameLen])

	offset += HostnameLen

	Message.ClientID = string(b[offset : offset+ClientIDLen])

	offset += ClientIDLen

	Message.ClientName = string(b[offset : offset+ClientNameLen])

	offset += ClientNameLen

	Message.Content = string(b[offset:])

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

func (m *Message) GetClientURL() string {
	return fmt.Sprintf("https://steamcommunity.com/profiles/%s", m.ClientID)
}

func (m *Message) GetClientColor() int {
	c := []byte(m.ClientID)

	i, _ := strconv.ParseInt(string(c[len(c)-6:]), 16, 64)

	return int(i)
}
