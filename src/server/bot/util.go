package bot

import (
	"regexp"

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
