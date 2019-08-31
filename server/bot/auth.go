package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Auth(fn exrouter.HandlerFunc) exrouter.HandlerFunc {
	return func(ctx *exrouter.Context) {
		guild, err := MessageGuild(ctx, ctx.Msg)

		if err != nil {
			logrus.WithField("error", err).Warn()
			ctx.Reply("Could not fetch guild: ", err)
		}

		member, err := ctx.Member(guild.ID, ctx.Msg.Author.ID)

		if err != nil {
			ctx.Reply("Could not fetch member: ", err)
		}

		if GuildMemberPermissions(member, guild)&discordgo.PermissionManageServer != 0 {
			ctx.Set("member", member)
			ctx.Set("guild", guild)

			fn(ctx)

			return
		}

		ctx.Reply("You do not have permission to use this command")
	}
}
