package relay

import (
	"database/sql"
	"fmt"
	"net"
	"sync"

	"github.com/rumblefrog/source-chat-relay/server/entity"
	"github.com/rumblefrog/source-chat-relay/server/filter"
	"github.com/rumblefrog/source-chat-relay/server/packet"

	"github.com/rumblefrog/source-chat-relay/server/protocol"
	"github.com/sirupsen/logrus"
)

var Instance *Relay

type Relay struct {
	clientMu sync.RWMutex

	Clients    map[*RelayClient]bool
	Router     chan protocol.Deliverable
	Bot        chan protocol.Deliverable
	Listener   net.Listener
	Statistics RelayStats
	Closed     bool
}

type RelayClient struct {
	Socket     net.Conn
	Data       chan []byte
	ID         string
	EntityName string
	Statistics RelayStats
}

type RelayStats struct {
	Incoming RelayTrafficStats
	Outgoing RelayTrafficStats
}

type RelayTrafficStats struct {
	MessageCount int
	ByteCount    int
}

func NewRelay() *Relay {
	return &Relay{
		Clients: make(map[*RelayClient]bool),
		Router:  make(chan protocol.Deliverable),
		Bot:     make(chan protocol.Deliverable),
	}
}

func (r *Relay) Listen(port int) error {
	var err error

	r.Listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	go r.StartRouting()
	go r.ProcessConnections()

	return nil
}

func (r *Relay) StartRouting() {
	for {
		if r.Closed {
			return
		}

		select {
		case message := <-r.Router:
			if filter.IsInFilter(message.Content()) {
				continue
			}

			r.clientMu.Lock()

			// Iterate connected clients
			for client := range r.Clients {
				tEntity, err := entity.GetEntity(client.ID)

				if err != nil {
					continue
				}

				if client.ID != message.Author() &&
					tEntity.CanReceiveType(message.Type()) &&
					tEntity.ReceiveIntersectsWith(entity.DeliverableSendChannels(message)) {
					select {
					case client.Data <- message.Marshal():
					default:
						close(client.Data)
						delete(r.Clients, client)
					}
				}
			}

			r.clientMu.Unlock()

			// Push to bot channel and it'll iterate Discord channels
			r.Bot <- message
		}
	}
}

func (r *Relay) ProcessConnections() {
	for {
		conn, err := r.Listener.Accept()

		if err != nil {
			if r.Closed {
				return
			}

			logrus.WithField("error", err).Warn("Unable to accept connection")

			return
		}

		logrus.WithField("address", conn.RemoteAddr()).Info("A client connected")

		client := &RelayClient{
			Socket: conn,
			Data:   make(chan []byte),
		}

		r.AddClient(client)

		go r.ListenClientReceive(client)
		go r.ListenClientSend(client)
	}
}

func (r *Relay) ListenClientReceive(c *RelayClient) {
	for {
		buffer := make([]byte, protocol.MAX_BUFFER_LENGTH)

		length, err := c.Socket.Read(buffer)

		if err != nil {
			r.RemoveClient(c)
			c.Socket.Close()
			break
		}

		if length > 0 {
			buffer = buffer[:length]

			r.Statistics.Incoming.ByteCount += length

			c.Statistics.Incoming.ByteCount += length

			r.HandlePacket(c, buffer)
		}
	}
}

func (r *Relay) ListenClientSend(c *RelayClient) {
	defer c.Socket.Close()

	for {
		select {
		case message, ok := <-c.Data:
			if !ok {
				// Exit for loop, execute the defer
				return
			}

			b, _ := c.Socket.Write(message)

			r.Statistics.Outgoing.ByteCount += b
			r.Statistics.Outgoing.MessageCount++

			c.Statistics.Outgoing.ByteCount += b
			c.Statistics.Outgoing.MessageCount++
		}
	}
}

func (r *Relay) HandlePacket(client *RelayClient, buffer []byte) {
	reader := packet.NewPacketReader(buffer)

	base, err := protocol.ParseBaseMessage(reader)

	if err != nil {
		return
	}

	r.Statistics.Incoming.MessageCount++
	client.Statistics.Incoming.MessageCount++

	if base.Type == protocol.MessageAuthenticate {
		authenticateMessage, err := protocol.ParseAuthenticateMessage(base, reader)
		authenticateResponseMessage := &protocol.AuthenticateMessageResponse{}

		if err != nil || len(authenticateMessage.Token) == 0 {
			authenticateResponseMessage.Response = protocol.AuthenticateDenied

			client.Socket.Write(authenticateResponseMessage.Marshal())

			logrus.WithField("address", client.Socket.RemoteAddr()).Warn("Client authentication failed")

			return
		}

		r.AuthenticateClient(client, authenticateMessage)

		authenticateResponseMessage.Response = protocol.AuthenticateSuccess

		client.Socket.Write(authenticateResponseMessage.Marshal())

		logrus.WithFields(logrus.Fields{
			"address":  client.Socket.RemoteAddr(),
			"hostname": client.EntityName,
			"id":       client.ID,
		}).Info("Client authenticated")

		return
	}

	// Switch case for everything else that requires auth prior

	if !client.Authenticated() {
		return
	}

	tEntity, err := entity.GetEntity(client.ID)

	if err != nil {
		return
	}

	if !tEntity.CanSendType(base.Type) {
		return
	}

	base.SenderID = client.ID

	// Reupdate entity name
	client.EntityName = base.EntityName

	switch base.Type {
	case protocol.MessageChat:
		msg, err := protocol.ParseChatMessage(base, reader)

		if err != nil {
			return
		}

		r.Router <- msg
	case protocol.MessageEvent:
		msg, err := protocol.ParseEventMessage(base, reader)

		if err != nil {
			return
		}

		r.Router <- msg
	default:
		// Malformed packet, we should not get anything else
		r.RemoveClient(client)
		client.Socket.Close()
	}
}

func (r *Relay) AddClient(c *RelayClient) {
	r.clientMu.Lock()
	defer r.clientMu.Unlock()

	r.Clients[c] = true
}

func (r *Relay) RemoveClient(c *RelayClient) {
	r.clientMu.Lock()
	defer r.clientMu.Unlock()

	if _, ok := r.Clients[c]; ok {
		close(c.Data)
		delete(r.Clients, c)
	}
}

func (r *Relay) AuthenticateClient(c *RelayClient, packet *protocol.AuthenticateMessage) {
	tEntity, err := entity.GetEntity(packet.Token)

	if err == sql.ErrNoRows {
		tEntity = &entity.Entity{
			ID: packet.Token,
		}

		if err = tEntity.Insert(); err != nil {
			logrus.WithField("error", err).Warn("Failed to create entity in database")
			return
		}
	} else if err != nil {
		logrus.WithField("error", err).Warn("Failed to fetch entity from database")
	}

	// Update database with new name upon auth
	tEntity.SetDisplayName(packet.BaseMessage.EntityName)

	c.ID = string(packet.Token)
	c.EntityName = packet.BaseMessage.EntityName
}
