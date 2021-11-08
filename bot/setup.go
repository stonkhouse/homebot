package bot

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

const (
	HOUSE_COLLECTION_PATH     = "houses"
	USER_COLLECTION_PATH      = "users"
	DBS_PAYLAH_LINK_SUBSTRING = "https://www.dbs.com.sg/personal/mobile/paylink/index.html?tranRef="
)

var (
	passwordChannel chan string
	paylahChannel   chan string
)

func (h *BotHandler) HandleStartHouse(m *telebot.Message) {
	var (
		members    []*User
		password   string
		paylahLink string
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
	_, _ = h.Bot.Send(m.Sender, "Hello! Thank you for using Homebot :)")
	_, _ = h.Bot.Send(m.Sender, "Please set your House Password")

	//	3. user sends password to the bot.
	h.Bot.Handle(telebot.OnText, h.setPassword)
	password = <-passwordChannel
	fmt.Printf("Password set: %s", password)

	//	4. bot prompts user to send paylah link
	paylahChannel = make(chan string)
	_, _ = h.Bot.Send(m.Sender, "Please enter your personal PayLah! Wallet Link")
	h.Bot.Handle(telebot.OnText, h.setPaylahLink)
	paylahLink = <-paylahChannel
	fmt.Printf("Password set: %s", paylahLink)

	//	5. add user to the db
	userObj := &User{
		ID:         m.Sender.ID,
		Username:   m.Sender.Username,
		PaylahLink: paylahLink,
	}
	fmt.Printf("User ID: %d", userObj.ID)
	fmt.Printf("Username: %s", userObj.Username)
	fmt.Printf("Paylah link: %s", userObj.PaylahLink)

	userDoc := h.Firestore.Collection(USER_COLLECTION_PATH).Doc(strconv.Itoa(m.Sender.ID))
	_, err := userDoc.Create(context.Background(), userObj)
	fmt.Printf("User Created!")
	if err != nil {
		fmt.Printf(err.Error())
	}
	members = append(members, userObj)

	houseObj := &House{
		HouseName: m.Chat.Title,
		Members:   members,
		Password:  password,
	}
	houseDoc := h.Firestore.Collection(HOUSE_COLLECTION_PATH).Doc(strconv.FormatInt(houseID, 10))
	_, err = houseDoc.Create(context.Background(), houseObj)
	fmt.Printf("House Created!")
	if err != nil {
		fmt.Printf(err.Error())
	}
	reply := fmt.Sprintf("Congratulations! House '%s' has been successfully created!\n"+
		"For users other than @%s, please type /join to join this house :D", m.Chat.Title, m.Sender.Username)
	_, _ = h.Bot.Send(m.Chat, reply)
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
	for {
		if m.Chat.Username == m.Sender.Username {
			passwordChannel <- m.Text
			_, _ = h.Bot.Reply(m, "Password successfully set!")
		}
	}
}
func (h *BotHandler) setPaylahLink(m *telebot.Message) {
	for {
		if strings.Contains(m.Text, DBS_PAYLAH_LINK_SUBSTRING) {
			paylahChannel <- m.Text
			_, _ = h.Bot.Reply(m, "PayLah Link successfully set!")
			//unset the handler
			h.Bot.Handle(telebot.OnText, func() {})
			return
		} else {
			_, _ = h.Bot.Reply(m, "Please enter a valid PayLah! link")
			return
		}
	}
}
