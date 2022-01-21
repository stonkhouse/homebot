package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
)

//TODO: Also allow for non recurring payments

const (
	PAYMENT_COLLECTION_PATH = "payments"
)

var (
	paymentChannel chan string
	senderUsername string
)

func (h *BotHandler) HandlePaymentReminder(m *telebot.Message) {
	senderUsername = m.Sender.Username
	paymentID := uuid.NewString()
	paymentChannel = make(chan string)
	_, _ = h.Bot.Reply(m, "Please enter the name of this monthly recurring payment")
	h.Bot.Handle(telebot.OnText, h.handleSetPaymentName)
	paymentName := <-paymentChannel
	_, _ = h.Bot.Send(m.Chat, "Please enter the date of the recurring payment (DD): ")
	h.Bot.Handle(telebot.OnText, h.handleSetPaymentDate)
	paymentDate := <-paymentChannel
	newPayment := &Payment{
		ID:      paymentID,
		Name:    paymentName,
		Date:    paymentDate,
		HouseID: m.Chat.ID,
	}
	//TODO: persist the payment on DB
	paymentDoc := h.Firestore.Collection(PAYMENT_COLLECTION_PATH).Doc(newPayment.ID)
	_, err := paymentDoc.Set(context.Background(), newPayment)
	if err != nil {
		fmt.Printf(err.Error())
	}
	message := fmt.Sprintf("Payment \"%s\" has been successfully created!\nIt's recurring on the %s of every month.", newPayment.Name, newPayment.Date)
	_, _ = h.Bot.Send(m.Chat, message)
}
func (h *BotHandler) handleSetPaymentName(m *telebot.Message) {
	for {
		if m.Sender.Username == senderUsername {
			paymentChannel <- m.Text
			_, _ = h.Bot.Reply(m, "Payment name set")
			//unset the handler
			h.Bot.Handle(telebot.OnText, func() {})
		}
		return
	}
}

func (h *BotHandler) handleSetPaymentDate(m *telebot.Message) {
	for {
		if m.Sender.Username == senderUsername {
			err := validateDate(m.Text)
			if err != nil {
				errMessage := fmt.Sprintf("Please enter a valid date: %s", err.Error())
				_, _ = h.Bot.Send(m.Chat, errMessage)
				return
			}
			paymentChannel <- m.Text
			_, _ = h.Bot.Reply(m, "Payment date set")
			//unset the handler
			h.Bot.Handle(telebot.OnText, func() {})
		}
		return
	}
}
func validateDate(date string) error {
	if len(date) != 2 {
		return errors.New("invalid date length")
	}
	dateVal, err := strconv.Atoi(date)
	if err != nil {
		return errors.New("date is not a number")
	} else {
		if dateVal > 31 || dateVal < 1 {
			return errors.New("date must be between 1 - 31")
		}
	}
	return nil
}
