package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotHandler struct {
	//menu     *tb.ReplyMarkup
	Bot *tb.Bot
}
type SetupOption struct {
	Option      string
	Description string
}
