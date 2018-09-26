package bot

func (b *DiscordBot) Listen() {
	for {
		select {
		case message := <-b.Data:
			for _, c := range b.RelayChannels {
				if c.CanReceive(message.Header.Sender.SendChannels) {
					// Yes
				}
			}
		}
	}
}

func (channel *RelayChannel) CanReceive(channels []int) bool {
	for c := range channel.ReceiveChannels {
		for c1 := range channels {
			if c == c1 {
				return true
			}
		}
	}

	return false
}
