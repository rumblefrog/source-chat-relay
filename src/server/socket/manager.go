package socket

import (
	"fmt"
	"net"

	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Router     chan []protocol.Message
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Socket   net.Conn
	Data     chan []byte
	Channels []int
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
			// case message := <-manager.Router:
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

			Header := protocol.NewHeader(message)

			Header.RequestLength = length

			switch Header.GetOPCode() {
			case protocol.MessageFrame:
				{
					go protocol.HandleMessage(message, Header)
				}
			case protocol.PingFrame:
				{
					go protocol.HandlePing(Header)
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

func (client *Client) HasChannel(channels []int) bool {
	for c := range client.Channels {
		for c1 := range channels {
			if c == c1 {
				return true
			}
		}
	}

	return false
}
