package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/protocol"

	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
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
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			for _, e := range repoEntity.GetEntities(repoEntity.All) {
				if e.Type == repoEntity.Channel && e.CanReceive(message.GetSendChannels()) {
					b.Session.ChannelMessageSendEmbed(e.ID, embed)
				}
			}
		}
	}
}
