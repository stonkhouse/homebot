package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
)

//TODO: Also allow for non recurring payments

const (
	PAYMENT_COLLECTION_PATH = "payments"
)

var (
	paymentChannel chan string
	amountChannel  chan int
	dateChannel    chan int
	senderUsername string
)

func (h *BotHandler) RegisterPaymentReminder() {
	fmt.Printf("\nSetting up payment reminder\n")
	//1. Fetch the list of payments to remind
	h.Cron = cron.New()
	payments := h.queryPayments()

	//2. For each payment, add function to send message to the group
	// This handles payments that are already in the DB before the bot started
	for _, payment := range payments {
		fmt.Printf("Payment: %s", payment.Name)
		h.createPaymentCron(h.Cron, payment)
	}
	h.Cron.Start()
}
func (h *BotHandler) createPaymentCron(c *cron.Cron, payment *Payment) {
	//TODO: Handle edge cases: start of the month
	spec := fmt.Sprintf("39 1 %d * *", payment.Date-1)
	// ! Note the implementation here
	go func(payment *Payment) {
		fmt.Printf("New CRON created for payment: %s\n", payment.Name)
		_, err := c.AddFunc(spec, func() { h.sendPaymentReminder(payment) })
		if err != nil {
			fmt.Printf("Cron error: %s", err.Error())
			return
		}
	}(payment)

}
func (h *BotHandler) sendPaymentReminder(p *Payment) {
	houseID := p.HouseID
	group := telebot.ChatID(houseID)
	message := fmt.Sprintf("Heyo! You have a bill due tomorrow: %s\nAmount: %d", p.Name, p.Amount)
	_, _ = h.Bot.Send(group, message)
}

func (h *BotHandler) HandleCreatePayment(m *telebot.Message) {
	var (
		paymentName   string
		paymentAmount int
		paymentDate   int
	)
	senderUsername = m.Sender.Username
	paymentID := uuid.NewString()

	paymentChannel = make(chan string)
	_, _ = h.Bot.Reply(m, "Please enter the name of this monthly recurring payment")
	h.Bot.Handle(telebot.OnText, h.handleSetPaymentName)
	paymentName = <-paymentChannel
	defer close(paymentChannel)

	amountChannel = make(chan int)
	_, _ = h.Bot.Reply(m, "Please enter the amount of this recurring payment")
	h.Bot.Handle(telebot.OnText, h.handleSetPaymentAmount)
	paymentAmount = <-amountChannel
	defer close(amountChannel)

	dateChannel = make(chan int)
	_, _ = h.Bot.Reply(m, "Please enter the date of the recurring payment (DD): ")
	h.Bot.Handle(telebot.OnText, h.handleSetPaymentDate)
	paymentDate = <-dateChannel
	defer close(dateChannel)

	newPayment := &Payment{
		ID:      paymentID,
		Name:    paymentName,
		Date:    paymentDate,
		Amount:  paymentAmount,
		HouseID: m.Chat.ID,
	}
	h.createPaymentCron(h.Cron, newPayment)
	paymentDoc := h.Firestore.Collection(PAYMENT_COLLECTION_PATH).Doc(newPayment.ID)
	_, err := paymentDoc.Set(context.Background(), newPayment)
	if err != nil {
		fmt.Printf(err.Error())
	}
	message := fmt.Sprintf("Payment \"%s\" has been successfully created!\nIt's recurring on %d every month.", newPayment.Name, newPayment.Date)
	_, _ = h.Bot.Send(m.Chat, message)
}
func (h *BotHandler) handleSetPaymentName(m *telebot.Message) {
	for {
		if m.Sender.Username == senderUsername {
			paymentChannel <- m.Text
			_, _ = h.Bot.Reply(m, "Payment name set")
			h.Bot.Handle(telebot.OnText, func() {})
			return
		}
	}
}
func (h *BotHandler) handleSetPaymentAmount(m *telebot.Message) {
	for {
		if m.Sender.Username == senderUsername {
			val, err := strconv.Atoi(m.Text)
			if err != nil {
				_, _ = h.Bot.Reply(m, "Please enter a valid amount!")
			}
			amountChannel <- val
			_, _ = h.Bot.Reply(m, "Payment amount set")
			h.Bot.Handle(telebot.OnText, func() {})
			return
		}
	}
}

func (h *BotHandler) handleSetPaymentDate(m *telebot.Message) {
	for {
		if m.Sender.Username == senderUsername {
			date, err := validateDate(m.Text)
			if err != nil {
				errMessage := fmt.Sprintf("Please enter a valid date: %s", err.Error())
				_, _ = h.Bot.Send(m.Chat, errMessage)
				return
			}
			dateChannel <- date
			_, _ = h.Bot.Reply(m, "Payment date set")
			h.Bot.Handle(telebot.OnText, func() {})
			return
		}
	}
}
func validateDate(date string) (int, error) {
	if len(date) != 2 {
		return 0, errors.New("invalid date length")
	}
	dateVal, err := strconv.Atoi(date)
	if err != nil {
		return 0, errors.New("date is not a number")
	} else {
		if dateVal > 31 || dateVal < 1 {
			return 0, errors.New("date must be between 1 - 31")
		}
	}
	return dateVal, nil
}

func (h *BotHandler) queryPayments() []*Payment {
	var payments []*Payment
	docRefs, err := h.Firestore.Collection(PAYMENT_COLLECTION_PATH).DocumentRefs(context.Background()).GetAll()
	if err != nil {
		fmt.Printf("Error in fetching payments: %s", err.Error())
		return nil
	}
	for _, docRef := range docRefs {
		paymentDoc, err := docRef.Get(context.Background())
		if err != nil {
			fmt.Printf("Error in getting payment Doc: %s", err.Error())
			return nil
		}
		payment := &Payment{}
		err = paymentDoc.DataTo(payment)
		if err != nil {
			fmt.Printf("Error in casting document to payment: %s", err.Error())
		}
		payments = append(payments, payment)
	}
	return payments
}
