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

			for _, e := range b.Cache {
				if e.CanReceive(message.GetSendChannels()) {
					b.Session.ChannelMessageSendEmbed(e.ID, embed)
				}
			}
		}
	}
}
