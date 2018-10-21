package bot

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

var (
	ChannelRegex = regexp.MustCompile("(?:<#)?([0-9]+)>?")
	UserRegex    = regexp.MustCompile("(?:<@!?)?([0-9]+)>?")
	RoleRegex    = regexp.MustCompile("(?:<@&)?([0-9]+)>?")
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

func CapitalChannelName(c *discordgo.Channel) string {
	nameBytes := []byte(c.Name)

	return string(bytes.ToUpper(nameBytes[:1])) + string(nameBytes[1:])
}

func ParseChannel(arg string) (string, bool) {
	if ChannelRegex.Match([]byte(arg)) {
		return ChannelRegex.FindStringSubmatch(arg)[1], true
	}

	return "", false
}

func ParseUser(arg string) (string, bool) {
	if UserRegex.Match([]byte(arg)) {
		return UserRegex.FindStringSubmatch(arg)[1], true
	}

	return "", false
}

func ParseRole(arg string) (string, bool) {
	if RoleRegex.Match([]byte(arg)) {
		return RoleRegex.FindStringSubmatch(arg)[1], true
	}

	return "", false
}

func TransformMentions(session *discordgo.Session, cid string, body string) string {
	channelID, found := ParseChannel(body)

	if found {
		channel, err := session.Channel(channelID)

		if err == nil {
			body = ChannelRegex.ReplaceAllString(body, fmt.Sprintf("#%s", channel.Name))
		}
	}

	userID, found := ParseUser(body)

	if found {
		user, err := session.User(userID)

		if err == nil {
			body = ChannelRegex.ReplaceAllString(body, fmt.Sprintf("@%s", user.Username))
		}
	}

	roleID, found := ParseRole(body)

	if found && session.StateEnabled {
		channel, err := session.Channel(cid)

		if err == nil {
			role, err := session.State.Role(channel.ID, roleID)

			if err == nil {
				body = ChannelRegex.ReplaceAllString(body, fmt.Sprintf("@%s", role.Name))
			}
		}
	}

	return body
}
