package net

// Message - payload consist of chat content
type Message struct {
	Hostname       string
	ClientSteamID3 string
	ClientName     string
}

// ParseMessage - Parses packet and returns a pointer to a Message
func ParseMessage( /* packet */ ) *Message {
	// Data
	return &Message{}
}
