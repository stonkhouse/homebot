package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

func RegisterBot(homebot *telebot.Bot) {
	fmt.Printf("Bot is starting: %s", homebot.Token)
	homebot.Handle("/start", func(m *telebot.Message) {
		homebot.Send(m.Sender, "HELLO")
	})
	homebot.Start()
}
