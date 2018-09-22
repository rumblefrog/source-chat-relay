package protocol

import (
	"fmt"
	"net"
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
				if connection.CanReceive(message.Header.Sender.SendChannels) {
					connection.Data <- []byte(message.ToString())
				}
			}
		}
	}
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
