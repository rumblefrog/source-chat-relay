package protocol

import (
	"database/sql"
	"fmt"
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
	Socket          net.Conn
	Data            chan []byte
	Token           string
	SendChannels    []int
	ReceiveChannels []int
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
		}
	}
}

func (manager *ClientManager) RegisterClient(client *Client, token []byte) {
	querystmt, err := database.DBConnection.Prepare("SELECT `receive_channels`, `send_channels` FROM `relay_entities` WHERE `id` = ? AND `type` = 0")

	if err != nil {
		log.Fatal("Failed to prepare query to register client")
		return
	}

	var (
		receiveChannels string
		sendChannels    string
	)

	err = querystmt.QueryRow(string(token)).Scan(&receiveChannels, &sendChannels)

	if err == sql.ErrNoRows {
		if _, err := database.DBConnection.Exec("INSERT INTO `relay_entities` (`id`) VALUES (?)", string(token)); err != nil {
			log.Fatal("Failed to create client in database", err)
		}
		return
	} else if err != nil {
		log.Fatal("Failed to query to register client", err)
		return
	}

	client.Token = string(token)

	client.ReceiveChannels = database.ParseChannels(receiveChannels)

	client.SendChannels = database.ParseChannels(sendChannels)

	log.Println(client.ReceiveChannels)
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
			fmt.Println("RECEIVED: " + string(message))

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
	for c := range client.ReceiveChannels {
		for c1 := range channels {
			if c == c1 {
				return true
			}
		}
	}

	return false
}
