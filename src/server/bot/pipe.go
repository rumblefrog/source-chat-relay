package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

func (b *DiscordBot) Listen() {
	for {
		select {
		case message := <-protocol.NetManager.Bot:
			embed := &discordgo.MessageEmbed{
				URL:         message.GetClientURL(),
				Title:       message.ClientName,
				Color:       message.GetClientColor(),
				Description: message.Content,
			}

			for _, c := range b.RelayChannels {
				if c.CanReceive(message.Header.Sender.SendChannels) {
					b.Session.ChannelMessageSendEmbed(c.ChannelID, embed)
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
