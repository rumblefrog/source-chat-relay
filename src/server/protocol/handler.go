package protocol

import (
	"log"
	"net"
)

func HandlePacket(conn net.Conn) {
	buffer := make([]byte, 2048)

	reqLen, err := conn.Read(buffer)

	if err != nil {
		log.Println("Error reading packet", err)
		return
	}

	Header := NewHeader(buffer)

	log.Println("DataLen", reqLen)
	log.Println("Data", Header)
}
