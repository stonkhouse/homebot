package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

func RegisterBot(homebot *telebot.Bot) {
	fmt.Printf("Bot is starting...")
	homebot.Handle("/start", func(m *telebot.Message) {
		_, err := homebot.Send(m.Sender, "HELLO")
		if err != nil {
			fmt.Printf("Error sending hello: %s", err)
			return
		}
	})
	homebot.Start()
}
