# dgrouter

Router to simplify command routing in discord bots.
[exrouter](https://github.com/Necroforger/dgrouter/tree/master/exrouter) provides some extra features like wrapping the router type to add command handlers which typecast to its own Context type.

If you make any interesting changes feel free to submit a Pull Request

- [exrouter godoc reference](https://godoc.org/github.com/Necroforger/dgrouter/exrouter)
- [exmiddleware godoc reference](https://godoc.org/github.com/Necroforger/dgrouter/exmiddleware)
- [dgrouter godoc reference](https://godoc.org/github.com/Necroforger/dgrouter)

## Features

- [Subroutes](https://github.com/Necroforger/dgrouter/blob/master/examples/subrouters/subrouters.go#L28)
- [Route grouping](https://github.com/Necroforger/dgrouter/blob/master/examples/middleware/middleware.go#L69)
- [Route aliases](https://github.com/Necroforger/dgrouter/blob/master/examples/soundboard/soundboard.go#L97)
- [Middleware](https://github.com/Necroforger/dgrouter/blob/master/examples/middleware/middleware.go#L38)
- [Regex matching](https://github.com/Necroforger/dgrouter/blob/master/examples/pingpong/pingpong.go#L39)

## example
```go 
router.On("ping", func(ctx *exrouter.Context) { ctx.Reply("pong")}).Desc("responds with pong")

router.On("avatar", func(ctx *exrouter.Context) {
	ctx.Reply(ctx.Msg.Author.AvatarURL("2048"))
}).Desc("returns the user's avatar")

router.Default = router.On("help", func(ctx *exrouter.Context) {
	var text = ""
	for _, v := range router.Routes {
		text += v.Name + " : \t" + v.Description + "\n"
	}
	ctx.Reply("```" + text + "```")
}).Desc("prints this help menu")
```
