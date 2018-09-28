package bot

import (
	"regexp"

	"github.com/rumblefrog/source-chat-relay/src/server/database"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

var ChannelRegex = regexp.MustCompile("^(?:<#)?([0-9]+)>?$")

func GuildMemberPermissions(member *discordgo.Member, guild *discordgo.Guild) (apermissions int) {
	if member.User.ID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			apermissions |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				apermissions |= role.Permissions
				break
			}
		}
	}

	if apermissions&discordgo.PermissionAdministrator != 0 {
		apermissions |= discordgo.PermissionAllChannel
	}

	return
}

func GetMessageGuild(c *exrouter.Context, m *discordgo.Message) (*discordgo.Guild, error) {
	channel, err := c.Channel(m.ChannelID)

	if err != nil {
		return nil, err
	}

	guild, err := c.Guild(channel.GuildID)

	if err != nil {
		return nil, err
	}

	return guild, nil
}

func (b *DiscordBot) GetRelayChannel(channelID string) *RelayChannel {
	for _, c := range b.RelayChannels {
		if c.ChannelID == channelID {
			return c
		}
	}

	return nil
}

func ParseChannel(arg string) (string, bool) {
	if ChannelRegex.Match([]byte(arg)) {
		return ChannelRegex.FindAllStringSubmatch(arg, -1)[0][1], true
	}

	return "", false
}

func DisplayEntity(ctx *exrouter.Context, entity *database.Entity, title string) {
	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
		Title: title,
		Color: 14795100,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Entity",
				Value: entity.ID,
			},
			&discordgo.MessageEmbedField{
				Name:  "Entity Type",
				Value: entity.Type.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "Receive Channels",
				Value: database.EncodeChannelsSep(entity.ReceiveChannels, ", "),
			},
			&discordgo.MessageEmbedField{
				Name:  "Send Channels",
				Value: database.EncodeChannelsSep(entity.SendChannels, ", "),
			},
			&discordgo.MessageEmbedField{
				Name:  "Created At",
				Value: entity.CreatedAt.String(),
			},
		},
	})
}
