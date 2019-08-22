package bot

import (
	"regexp"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

var (
	ChannelRegexExplicit = regexp.MustCompile("^(?:<#)?([0-9]+)>?$")
)

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

func ParseChannel(arg string) (string, bool) {
	if ChannelRegexExplicit.Match([]byte(arg)) {
		return ChannelRegexExplicit.FindStringSubmatch(arg)[1], true
	}

	return "", false
}
