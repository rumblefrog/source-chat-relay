package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

func Auth(fn exrouter.HandlerFunc) exrouter.HandlerFunc {
	return func(ctx *exrouter.Context) {
		member, err := ctx.Member(ctx.Msg.GuildID, ctx.Msg.Author.ID)

		if err != nil {
			ctx.Reply("Could not fetch member: ", err)
		}

		guild, err := ctx.Guild(ctx.Msg.GuildID)

		if err != nil {
			ctx.Reply("Could not fetch guild: ", err)
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
