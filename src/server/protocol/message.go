package protocol

const (
	HostnameLen   = 64
	ClientIDLen   = 64
	ClientNameLen = 32
)

type Message struct {
	Hostname   string
	ClientID   string
	ClientName string
}

func ParseMessage(b []byte) *Message {
	var (
		hostname   = make([]byte, HostnameLen)
		clientid   = make([]byte, ClientIDLen)
		clientname = make([]byte, ClientNameLen)
		offset     = 0
	)

	for i := 0; i < HostnameLen; i++ {
		hostname = append(hostname, b[offset])
		offset++
	}

	for i := 0; i < ClientIDLen; i++ {
		clientid = append(clientid, b[offset])
		offset++
	}

	for i := 0; i < ClientNameLen; i++ {
		clientname = append(clientname, b[offset])
		offset++
	}

	return &Message{
		Hostname:   string(hostname),
		ClientID:   string(clientid),
		ClientName: string(clientname),
	}
}

func (m *Message) GetHostname() string {
	return m.Hostname
}

func (m *Message) GetClientID() string {
	return m.ClientID
}

func (m *Message) GetClientName() string {
	return m.ClientName
}
