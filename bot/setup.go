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
	HOUSE_COLLECTION_PATH  = "houses"
	USER_COLLECTION_PATH   = "users"
	REVOLUT_LINK_SUBSTRING = "https://revolut.me/"
)

var (
	passwordChannel chan string
	revolutChannel  chan string
)

func (h *BotHandler) HandleStartHouse(m *telebot.Message) {
	var (
		members     []*User
		password    string
		revolutLink string
	)
	// 0. If chat is not a group, exit
	if !m.FromGroup() {
		_, _ = h.Bot.Reply(m, "You can only run this command within a group!")
		return
	}

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

	//TODO: Extract the set Revolut link to a function to avoid repetition

	//	4. bot prompts user to send Revolut link
	revolutChannel = make(chan string)
	_, _ = h.Bot.Send(m.Sender, "Please enter your personal Revolut link")
	h.Bot.Handle(telebot.OnText, h.setRevolutLink)
	revolutLink = <-revolutChannel

	//	5. add user to the db
	userObj := &User{
		ID:          m.Sender.ID,
		Username:    m.Sender.Username,
		RevolutLink: revolutLink,
		HouseID:     houseID,
	}
	fmt.Printf("User house ID: %d", userObj.HouseID)
	fmt.Printf("Username: %s", userObj.Username)
	fmt.Printf("Paylah link: %s", userObj.RevolutLink)

	userDoc := h.Firestore.Collection(USER_COLLECTION_PATH).Doc(strconv.Itoa(m.Sender.ID))
	_, err := userDoc.Set(context.Background(), userObj)
	fmt.Printf("User Created!")
	if err != nil {
		fmt.Printf(err.Error())
	}
	members = append(members, userObj)

	houseObj := &House{
		ID:        houseID,
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
	//TODO: Extract all fromgroup check to a middleware
	if !m.FromGroup() {
		_, _ = h.Bot.Reply(m, "You can only run this command within a group!")
		return
	}
	userID := m.Sender.ID
	user := h.queryUserByID(userID)
	house := h.queryHouseByID(m.Chat.ID)

	//0a. House should exist
	if house == nil {
		_, _ = h.Bot.Reply(m, "House doesn't exist. Please run the /start_house command first :)")
		return
	}
	//0b. User should either not exist or exist but not bound to a house
	if user != nil {
		if user.HouseID != 0 {
			userHouse := h.queryHouseByID(user.HouseID)
			reply := fmt.Sprintf("You already belong to House: %s\n Please type /leave in that house before joining a new house!", userHouse.HouseName)
			_, _ = h.Bot.Reply(m, reply)
			return
		}
	}

	//1. Bot PM user "pls enter house password"
	msg := fmt.Sprintf("Welcome to House: %s\nPlease enter the house password!", house.HouseName)
	_, _ = h.Bot.Send(m.Sender, msg)
	houseID := m.Chat.ID
	houseObj := h.queryHouseByID(houseID)
	if houseObj == nil {
		return
	} else {
		passwordChannel = make(chan string)
		h.Bot.Handle(telebot.OnText, h.checkPassword)
		password := <-passwordChannel

		//2. if successful, request for paylah link
		if password == house.Password {
			revolutChannel = make(chan string)
			_, _ = h.Bot.Send(m.Sender, "Please enter your personal Revolut Llink")
			h.Bot.Handle(telebot.OnText, h.setRevolutLink)
			paylahLink := <-revolutChannel

			userObj := &User{
				ID:          m.Sender.ID,
				Username:    m.Sender.Username,
				RevolutLink: paylahLink,
				HouseID:     houseID,
			}
			fmt.Printf("User house ID: %d", userObj.HouseID)
			fmt.Printf("Username: %s", userObj.Username)
			fmt.Printf("Paylah link: %s", userObj.RevolutLink)

			userDoc := h.Firestore.Collection(USER_COLLECTION_PATH).Doc(strconv.Itoa(m.Sender.ID))
			_, err := userDoc.Set(context.Background(), userObj)
			fmt.Printf("User Created!")
			if err != nil {
				fmt.Printf(err.Error())
			}
			house.Members = append(house.Members, userObj)
			houseDoc := h.Firestore.Collection(HOUSE_COLLECTION_PATH).Doc(strconv.FormatInt(house.ID, 10))
			_, err = houseDoc.Set(context.Background(), house)
			if err != nil {
				reply := fmt.Sprintf("Error in updating members: %s", err.Error())
				_, _ = h.Bot.Reply(m, reply)
			}
			fmt.Printf("User added to the house")
			reply := fmt.Sprintf("User @%s has been successfully added to this house.\nWelcome!????", userObj.Username)
			_, _ = h.Bot.Reply(m, reply)
		} else {
			_, _ = h.Bot.Send(m.Sender, "Sorry, the password that you entered is not correct, please type /join in your house group again :(")
		}

	}
}
func (h *BotHandler) HandleOnAddToGroup(m *telebot.Message) {
	houseID := m.Chat.ID
	//	1. Check the database, see if group ID exists in the DB
	house := h.queryHouseByID(houseID)
	//	1a. if group exists, bot: "Hello welcome back", move to 3
	if house != nil {
		_, _ = h.Bot.Send(m.Chat, "Heyooo it me again :D")
	} else {
		//	1b. if group doesn't exist, bot: "Hello thanks for using homebot"
		_, _ = h.Bot.Send(m.Chat, "Heyoooo thanks for using Homebot! ????\nSlide into my DM and run /start.\nThen, run /start_house to start house ????")
	}
}
func (h *BotHandler) HandleLeave(m *telebot.Message) {
	if !m.FromGroup() {
		_, _ = h.Bot.Reply(m, "You can only run this command within a group!")
		return
	}
	userID := m.Sender.ID
	houseID := m.Chat.ID
	user := h.queryUserByID(userID)
	house := h.queryHouseByID(houseID)
	houseEmpty := false
	if user == nil {
		_, _ = h.Bot.Reply(m, "You don't exist yet ????")
		return
	}
	switch user.HouseID {
	case 0:
		_, _ = h.Bot.Reply(m, "You don't belong to any house!")
		return

	case houseID:
		updatedMembers := h.removeUserFromHouse(userID, house.Members)
		house.Members = updatedMembers
		houseDoc := h.Firestore.Collection(HOUSE_COLLECTION_PATH).Doc(strconv.FormatInt(houseID, 10))

		//Remove house when it's empty
		if len(house.Members) == 0 {
			_, err := houseDoc.Delete(context.Background())
			houseEmpty = true
			if err != nil {
				reply := fmt.Sprintf("Error in deleting house: %s", err.Error())
				_, _ = h.Bot.Reply(m, reply)
				return
			}
		} else {
			//else, update it
			_, err := houseDoc.Set(context.Background(), house)
			if err != nil {
				reply := fmt.Sprintf("Error in updating members: %s", err.Error())
				_, _ = h.Bot.Reply(m, reply)
				return
			}
		}

		user.HouseID = 0
		userDoc := h.Firestore.Collection(USER_COLLECTION_PATH).Doc(strconv.Itoa(userID))
		_, err := userDoc.Delete(context.Background())
		if err != nil {
			reply := fmt.Sprintf("Error in deleting user: %s", err.Error())
			_, _ = h.Bot.Reply(m, reply)
		}
		_, _ = h.Bot.Reply(m, "You have been successfully removed from this house :\")")
		if houseEmpty {
			_, _ = h.Bot.Send(m.Chat, "Since there are no more members in the house, this house will be deleted!")
		}
		return

	//default: If user.houseID != houseID
	default:
		house := h.queryHouseByID(user.HouseID)
		reply := fmt.Sprintf("You belong to House: %s, please run /leave in the right house ????", house.HouseName)
		_, _ = h.Bot.Reply(m, reply)
		return
	}
}
func (h *BotHandler) queryUserByID(userID int) *User {
	userDoc, err := h.Firestore.Doc(USER_COLLECTION_PATH + "/" + strconv.Itoa(userID)).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			fmt.Printf("User not found")
		} else {
			fmt.Printf("Error in fetching user: %s\n", err)
		}
		return nil
	}
	user := &User{}
	err = userDoc.DataTo(user)
	if err != nil {
		fmt.Print("Error in storing data to house object")
		return nil
	}
	return user
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
func (h *BotHandler) checkPassword(m *telebot.Message) {
	for {
		if m.Chat.Username == m.Sender.Username {
			passwordChannel <- m.Text
		}
	}
}
func (h *BotHandler) setRevolutLink(m *telebot.Message) {
	for {
		if m.Chat.Username == m.Sender.Username {
			if strings.Contains(m.Text, REVOLUT_LINK_SUBSTRING) {
				revolutChannel <- m.Text
				_, _ = h.Bot.Reply(m, "Revolut Link successfully set!")
				//unset the handler
				h.Bot.Handle(telebot.OnText, func() {})
				return
			} else {
				_, _ = h.Bot.Reply(m, "Please enter a valid Revolut link")
				return
			}
		}
	}
}

func (h *BotHandler) removeUserFromHouse(userID int, members []*User) []*User {
	j := 0
	for _, member := range members {
		if member.ID != userID {
			members[j] = member
			j++
		}
	}
	members = members[:j]
	return members
}
