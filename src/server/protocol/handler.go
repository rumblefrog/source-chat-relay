package protocol

import (
	"log"
)

func HandleMessage(b []byte, h *Header) {
	Message := ParseMessage(b, h)

	log.Println(Message.GetClientName())

	//TODO: Obtain channels & store in msg struct

	// log.Printf("Hostname: %s \n", Message.Hostname)
	// log.Printf("ID: %s \n", Message.ClientID)
	// log.Printf("Name: %s \n", Message.ClientName)
	// log.Printf("Content: %s", Message.Content)
}

func HandlePing(h *Header) {

}
