package bot

import (
	"cloud.google.com/go/firestore"
	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotHandler struct {
	//menu     *tb.ReplyMarkup
	Bot       *tb.Bot
	Firestore *firestore.Client
	Cron      *cron.Cron
}
type SetupOption struct {
	Option      string
	Text        string
	Description string
}

type House struct {
	ID        int64
	HouseName string
	Members   []*User
	Password  string
}

type User struct {
	ID         int
	Username   string
	PaylahLink string
	HouseID    int64
}

type Payment struct {
	ID      string
	Name    string
	Date    int
	Amount  int
	HouseID int64
	Payees  []User
}
