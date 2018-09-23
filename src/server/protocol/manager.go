package protocol

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/rumblefrog/source-chat-relay/src/server/database"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Router     chan *Message
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Socket          net.Conn
	Data            chan []byte
	SendChannels    []int
	ReceiveChannels []int
}

func (manager *ClientManager) Start() {
	for {
		select {
		case connection := <-manager.Register:
			manager.RegisterClient(connection)
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
				if connection.CanReceive(message.Header.Sender.SendChannels) {
					connection.Data <- []byte(message.ToString())
				}
			}
		}
	}
}

func (manager *ClientManager) RegisterClient(client *Client) {
	querystmt, err := database.DBConnection.Prepare("SELECT * FROM `relay_entities` WHERE `source` = ? AND `type` = 0")

	if err != nil {
		log.Panic("Failed to prepare query to register client")
		return
	}

	data := database.RelayEntities{}

	err = querystmt.QueryRow(client.Socket.RemoteAddr()).Scan(&data)

	if err == sql.ErrNoRows {
		insertstmt, err := database.DBConnection.Prepare("INSERT INTO `relay_entities` (`source`) VALUES (?)")

		if err != nil {
			log.Panic("Failed to prepare create client statement")
			return
		}

		_, err = insertstmt.Exec()

		if err != nil {
			log.Panic("Failed to create client in database")
			return
		}

		manager.Clients[client] = true

		return
	} else if err != nil {
		log.Panic("Failed to query to register client")
		return
	}

	client.ReceiveChannels = database.ParseChannels(data.ReceiveChannels)

	client.SendChannels = database.ParseChannels(data.SendChannels)

	manager.Clients[client] = true
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
			// TODO: Remove
			fmt.Println("RECEIVED: " + string(message))

			Header := NewHeader(message)

			Header.Sender = client

			Header.RequestLength = length

			switch Header.GetOPCode() {
			case MessageFrame:
				{
					go HandleMessage(message, Header)
				}
			case PingFrame:
				{
					go HandlePing(Header)
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
	for c := range client.ReceiveChannels {
		for c1 := range channels {
			if c == c1 {
				return true
			}
		}
	}

	return false
}
