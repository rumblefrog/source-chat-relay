package main

import (
	"./bot"
	"./helper"
	"./socket"
)

func main() {
	helper.LoadConfig()

	socket.InitSocket()

	bot.InitBot()
}
