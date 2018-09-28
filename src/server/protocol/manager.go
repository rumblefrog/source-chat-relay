package protocol

import (
	"database/sql"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/src/server/database"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Bot        chan *Message
	Router     chan *Message
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Socket net.Conn
	Data   chan []byte
	Entity *database.Entity
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
		case message := <-manager.Broadcast:
			for connection := range manager.Clients {
				select {
				case connection.Data <- message:
				default:
					close(connection.Data)
					delete(manager.Clients, connection)
				}
			}
		case message := <-manager.Router:
			for connection := range manager.Clients {
				if connection.CanReceive(message.GetSendChannels()) {
					connection.Data <- []byte(message.ToString())
					if message.Overwrite == nil {
						manager.Bot <- message
					}
				}
			}
			log.WithFields(log.Fields{
				"Hostname":    message.Hostname,
				"Client Name": message.ClientName,
				"Client ID":   message.ClientID,
			}).Debug()
		}
	}
}

func (manager *ClientManager) RegisterClient(client *Client, token []byte) {
	entity, err := database.FetchEntity(string(token))

	if err == sql.ErrNoRows {
		entity = &database.Entity{
			ID:   string(token),
			Type: database.Server,
		}

		if _, err = entity.CreateEntity(); err != nil {
			log.WithField("error", err).Warn("Failed to create entity in database")
			return
		}
	} else if err != nil {
		log.WithField("error", err).Warn("Failed to fetch entity from database")
	}

	client.Entity = entity
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

			// TODO: Remove
			log.WithField("Message", string(message)).Debug("Received Message")

			Header := NewHeader(message)

			Header.Sender = client

			Header.RequestLength = length

			switch Header.GetOPCode() {
			case AuthenticateFrame:
				{
					go manager.RegisterClient(client, message[1:])
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

func (client *Client) CanReceive(channels []int) bool {
	for c := range client.Entity.ReceiveChannels {
		for c1 := range channels {
			if c == c1 {
				return true
			}
		}
	}

	return false
}
