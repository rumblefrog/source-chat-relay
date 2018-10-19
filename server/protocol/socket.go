package protocol

import (
	"fmt"
	"net"

	"github.com/rumblefrog/source-chat-relay/server/config"
	log "github.com/sirupsen/logrus"
)

var (
	NetListener net.Listener
	NetManager  *ClientManager
)

func init() {
	var err error

	NetListener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.Conf.General.Port))

	if err != nil {
		log.WithField("error", err).Fatal("Unable to start socket server")
		return
	}

	NetManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Bot:        make(chan *Message),
		Router:     make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	go AcceptConnections()
}

func AcceptConnections() {
	go NetManager.Start()

	for {
		conn, err := NetListener.Accept()

		if err != nil {
			log.WithField("error", err).Warn("Unable to accept connection")
			return
		}

		log.WithField("address", conn.RemoteAddr()).Info("A client connected")

		client := &Client{
			Socket: conn,
			Data:   make(chan []byte),
		}

		NetManager.Register <- client

		go NetManager.Receive(client)

		go NetManager.Send(client)
	}
}
