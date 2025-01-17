package main

import (
	"github.com/d1mk9/tgPathToMeBot/configs"
	"github.com/d1mk9/tgPathToMeBot/internal/bot"
)

func main() {
	configs.LoadConfig()
	bot.StartBot()
}
