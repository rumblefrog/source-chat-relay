package relay

import "fmt"

func (c *RelayClient) Authenticated() bool {
	return len(c.ID) != 0
}

func (s RelayTrafficStats) String() string {
	return fmt.Sprintf("Messages: %d (%f MB)", s.MessageCount, float64(s.ByteCount)/(1024*1024))
}
