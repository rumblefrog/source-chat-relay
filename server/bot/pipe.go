package bot

import (
	"fmt"
	"time"

	"github.com/rumblefrog/source-chat-relay/server/config"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/protocol"

	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
)

func (b *DiscordBot) Listen() {
	for {
		select {
		case message := <-protocol.NetManager.Bot:

			var embed *discordgo.MessageEmbed

			if !config.Conf.Bot.SimpleMessage {
				embed = &discordgo.MessageEmbed{
					Color:       message.GetClientColor(),
					Description: message.Content,
					Timestamp:   time.Now().Format(time.RFC3339),
					Author: &discordgo.MessageEmbedAuthor{
						Name: message.ClientName,
						URL:  message.GetClientURL(),
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: message.Hostname,
					},
				}
			}

			for _, e := range repoEntity.GetEntities(repoEntity.All) {
				if e.Type == repoEntity.Channel && e.CanReceive(message.GetSendChannels()) {
					if !config.Conf.Bot.SimpleMessage {
						b.Session.ChannelMessageSendEmbed(e.ID, embed)
					} else {
						b.Session.ChannelMessageSend(e.ID, fmt.Sprintf("%s: %s", message.ClientName, message.Content))
					}
				}
			}
		}
	}
}
