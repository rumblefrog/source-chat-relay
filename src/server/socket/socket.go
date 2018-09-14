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

	NetListener, err = net.Listen("tcp", fmt.Sprintf(":%d", helper.Conf.Port))

	if err != nil {
		log.Panic("Unable to start socket server", err)
		return
	}

	go AcceptConnections()
}

func AcceptConnections() {
	for {
		conn, err := NetListener.Accept()

		if err != nil {
			log.Println("Unable to accept connection", err)
			return
		}

		go protocol.HandlePacket(conn)
	}
}
