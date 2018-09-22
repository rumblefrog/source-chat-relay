package socket

import (
	"fmt"
	"log"
	"net"

	"github.com/rumblefrog/source-chat-relay/src/server/helper"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

var NetListener net.Listener

func InitSocket() {
	var err error

	NetListener, err = net.Listen("tcp", fmt.Sprintf(":%d", helper.Conf.General.Port))

	if err != nil {
		log.Panic("Unable to start socket server", err)
		return
	}

	go AcceptConnections()
}

func AcceptConnections() {
	manager := ClientManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Router:     make(chan []protocol.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	go manager.Start()

	for {
		conn, err := NetListener.Accept()

		if err != nil {
			log.Println("Unable to accept connection", err)
			return
		}

		log.Println(fmt.Sprintf("%s connected", conn.RemoteAddr()))

		client := &Client{
			Socket: conn,
			Data:   make(chan []byte),
		}

		// TODO: Look up in data for the server's channels

		manager.Register <- client

		go manager.Receive(client)

		go manager.Send(client)

		go protocol.HandlePacket(conn)
	}
}
