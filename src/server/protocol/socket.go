package protocol

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/src/server/helper"
)

var (
	NetListener net.Listener
	NetManager  *ClientManager
)

func init() {
	var err error

	NetListener, err = net.Listen("tcp", fmt.Sprintf(":%d", helper.Conf.General.Port))

	if err != nil {
		log.Fatal("Unable to start socket server", err)
		return
	}

	go AcceptConnections()
}

func AcceptConnections() {
	NetManager := &ClientManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Router:     make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	go NetManager.Start()

	for {
		conn, err := NetListener.Accept()

		if err != nil {
			log.Warn("Unable to accept connection ", err)
			return
		}

		log.WithField("address", conn.RemoteAddr()).Info("A client connected")

		client := &Client{
			Socket: conn,
			Data:   make(chan []byte),
		}

		// TODO: Look up in data for the server's channels

		NetManager.Register <- client

		go NetManager.Receive(client)

		go NetManager.Send(client)
	}
}
