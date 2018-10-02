package protocol

import (
	"database/sql"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/server/filter"
	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Bot        chan *Message
	Router     chan *Message
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Socket net.Conn
	Data   chan []byte
	ID     string
}

func (manager *ClientManager) Start() {
	for {
		select {

		case connection := <-manager.Register:
			manager.Clients[connection] = true

		case connection := <-manager.Unregister:
			if _, ok := manager.Clients[connection]; ok {
				close(connection.Data)
				delete(manager.Clients, connection)
			}

		case message := <-manager.Router:
			if filter.IsInFilter(message.Content) {
				return
			}

			for connection := range manager.Clients {
				entity, err := repoEntity.GetEntity(connection.ID, repoEntity.Server)

				if err != nil {
					continue
				}

				if !connection.IsNotFromSelf(message) && entity.CanReceive(message.GetSendChannels()) {
					select {
					case connection.Data <- []byte(message.ToString()):
					default:
						close(connection.Data)
						delete(manager.Clients, connection)
					}
				}
			}

			if message.Overwrite == nil {
				manager.Bot <- message
			}
		}
	}
}

func (client *Client) Register(token []byte) {
	entity, err := repoEntity.GetEntity(string(token), repoEntity.Server)

	if err == sql.ErrNoRows {
		entity = &repoEntity.Entity{
			ID:   string(token),
			Type: repoEntity.Server,
		}

		if err = entity.Insert(); err != nil {
			log.WithField("error", err).Warn("Failed to create entity in database")
			return
		}
	} else if err != nil {
		log.WithField("error", err).Warn("Failed to fetch entity from database")
	}

	client.ID = string(token)
}

func (manager *ClientManager) Receive(client *Client) {
	for {
		message := make([]byte, 2048)
		length, err := client.Socket.Read(message)
		if err != nil {
			manager.Unregister <- client
			client.Socket.Close()
			break
		}
		if length > 0 {
			message = message[:length]

			log.WithField("message", string(message)).Debug("Received Message")

			Header := NewHeader(message)

			Header.Sender = client

			Header.RequestLength = length

			switch Header.GetOPCode() {
			case AuthenticateFrame:
				{
					go client.Register(message[1:])
				}
			case MessageFrame:
				{
					go manager.HandleMessage(message, Header)
				}
			}
		}
	}
}

func (manager *ClientManager) Send(client *Client) {
	defer client.Socket.Close()
	for {
		select {
		case message, ok := <-client.Data:
			if !ok {
				return
			}
			client.Socket.Write(message)
		}
	}
}

func (client *Client) IsNotFromSelf(message *Message) bool {
	if message.Overwrite != nil {
		return false
	}

	return client.ID == message.Header.Sender.ID
}
