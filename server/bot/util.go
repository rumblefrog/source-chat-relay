package bot

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

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

func TransformMentions(session *discordgo.Session, cid string, body string) string {
	if ChannelRegex.Match([]byte(body)) {
		matches := ChannelRegex.FindAllStringSubmatch(body, -1)

		n := len(matches)

		for i := 0; i < n; i++ {
			channel, err := session.Channel(matches[i][1])

			if err == nil {
				body = strings.Replace(body, matches[i][0], fmt.Sprintf("#%s", channel.Name), -1)
			}
		}
	}

	// Role match has to be before user, otherwise UserRegex will partial match role
	if RoleRegex.Match([]byte(body)) {
		channel, err := session.Channel(cid)

		if err == nil {
			matches := RoleRegex.FindAllStringSubmatch(body, -1)

			n := len(matches)

			for i := 0; i < n; i++ {
				role, err := session.State.Role(channel.GuildID, matches[i][1])

				if err == nil {
					body = strings.Replace(body, matches[i][0], fmt.Sprintf("@%s", role.Name), -1)
				}
			}
		}
	}

	if UserRegex.Match([]byte(body)) {
		matches := UserRegex.FindAllStringSubmatch(body, -1)

		n := len(matches)

		for i := 0; i < n; i++ {
			user, err := session.User(matches[i][1])

			if err == nil {
				body = strings.Replace(body, matches[i][0], fmt.Sprintf("@%s", user.Username), -1)
			}
		}
	}

	return body
}
