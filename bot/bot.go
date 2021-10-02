package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

func (h *BotHandler) RegisterBot() {
	fmt.Printf("Bot is starting...")
	fmt.Printf(h.Bot.Token)
	h.Bot.Handle("/start", h.HandleStart)
	h.Bot.Handle(telebot.OnQuery, h.HandleInlineQuery)
	h.Bot.Handle(telebot.OnChosenInlineResult, h.HandleSetup)
	h.Bot.Start()
}

func (h *BotHandler) HandleStart(m *telebot.Message) {
	_, err := h.Bot.Send(m.Sender, "Welcome to homebot! Type /help to see the list of available commands")
	if err != nil {
		fmt.Printf("Error sending hello: %s", err)
		return
	}
}
