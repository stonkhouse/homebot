package bot

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
)

const (
	HOUSE_COLLECTION_PATH = "houses"
)

var (
	passwordChannel chan string
)

func (h *BotHandler) HandleStartHouse(m *telebot.Message) {
	var (
		members  []*User
		password string
	)

	houseID := m.Chat.ID
	house := h.queryHouseByID(houseID)

	//	1. If house ID exists, bot replies "house already exists"
	if house != nil {
		reply := fmt.Sprintf("House '%s' already exists!", m.Chat.Title)
		_, _ = h.Bot.Reply(m, reply)
		return
	}
	//	2. Bot replies person to get One-Time-Password
	fmt.Println("Setting password...")
	passwordChannel = make(chan string)
	_, _ = h.Bot.Send(m.Sender, "Please set your House Password")

	//	3. user sends password to the bot.
	h.Bot.Handle(telebot.OnText, h.setPassword)
	password = <-passwordChannel
	//	4. bot prompts user to send paylah link
	fmt.Printf("Password set: %s", password)
	_, _ = h.Bot.Send(m.Sender, "Please enter your personal PayLah! Wallet Link")

	//	5. add the paylah link to the db

	houseDoc := h.Firestore.Collection(HOUSE_COLLECTION_PATH).Doc(strconv.FormatInt(houseID, 10))
	houseObj := &House{
		HouseName: m.Chat.Title,
		Members:   members,
		Password:  password,
	}
	_, err := houseDoc.Create(context.Background(), houseObj)
	fmt.Printf("House Created!")
	if err != nil {
		fmt.Printf("")
	}
	return
}
func (h *BotHandler) HandleJoin(m *telebot.Message) {
	//1. Bot PM user "pls enter one time password"
	//2. if successful, request for paylah link
	//3. add the paylah link to the DB
	//4.
}
func (h *BotHandler) HandleOnAddToGroup(m *telebot.Message) {
	//	1. Check the database, see if group ID exists in the DB
	//	1a. if group exists, bot: "Hello welcome back", move to 3
	//	1b. if group doesn't exist, bot: "Hello thanks for using homebot"
}

func (h *BotHandler) queryHouseByID(houseID int64) *House {
	houseDoc, err := h.Firestore.Doc(HOUSE_COLLECTION_PATH + "/" + strconv.FormatInt(houseID, 10)).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			fmt.Printf("House not found")
		} else {
			fmt.Printf("Error in fetching house: %s\n", err)
		}
		return nil
	}
	house := &House{}
	err = houseDoc.DataTo(house)
	if err != nil {
		fmt.Print("Error in storing data to house object")
		return nil
	}
	return house
}
func (h *BotHandler) setPassword(m *telebot.Message) {
	passwordChannel <- m.Text
	_, _ = h.Bot.Send(m.Sender, "Password successfully set!")
}
