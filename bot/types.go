package bot

import (
	"cloud.google.com/go/firestore"
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotHandler struct {
	//menu     *tb.ReplyMarkup
	Bot       *tb.Bot
	Firestore *firestore.Client
}
type SetupOption struct {
	Option      string
	Text        string
	Description string
}

type House struct {
	HouseName string
	Members   []*User
	Password  string
}

type User struct {
	ID         int
	Username   string
	PaylahLink string
}
