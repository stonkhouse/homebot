package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

func (h *BotHandler) RegisterBot() {
	fmt.Printf("Bot is starting...")
	fmt.Printf(h.Bot.Token)
	//Setup commands
	h.Bot.Handle("/start", h.HandleStart)
	h.Bot.Handle("/start_house", h.HandleStartHouse)
	h.Bot.Handle(telebot.OnAddedToGroup, h.HandleOnAddToGroup)
	h.Bot.Start()
}

func (h *BotHandler) HandleStart(m *telebot.Message) {
	_, err := h.Bot.Send(m.Chat, "Welcome to homebot! Type /help to see the list of available commands")
	if err != nil {
		fmt.Printf("Error sending hello: %s", err)
		return
	}
}
