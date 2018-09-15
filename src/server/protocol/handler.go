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

	Header.RequestLength = reqLen

	switch Header.GetOPCode() {
	case MessageFrame:
		{
			go HandleMessage(buffer, Header)
		}
	case PingFrame:
		{
			go HandlePing(Header)
		}
	}
}

func HandleMessage(b []byte, h *Header) {
	Message := ParseMessage(b, h)

	log.Println(Message.GetClientName())

	// log.Printf("Hostname: %s \n", Message.Hostname)
	// log.Printf("ID: %s \n", Message.ClientID)
	// log.Printf("Name: %s \n", Message.ClientName)
	// log.Printf("Content: %s", Message.Content)
}

func HandlePing(h *Header) {

}
