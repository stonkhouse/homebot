package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"gopkg.in/tucnak/telebot.v2"
	. "stonkhouse/stonkbot/bot"
	c "stonkhouse/stonkbot/config"
	"time"
)

func main() {
	var (
		config c.Configurations
	)

	//Reading configuration files
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error decoding config file: %s\n", err)
	}

	//initializing bot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.Telebot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	//initializing firebase
	opt := option.WithCredentialsFile(config.Firebase.ConfigPath)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("Error initializing Firebase App: %s\n", err)
	}
	firestore, err := firebaseApp.Firestore(context.Background())
	if err != nil {
		fmt.Printf("Error initializing Firestore: %s\n", err)
	}
	homebotHandler := &BotHandler{
		Bot:       bot,
		Firestore: firestore,
	}
	homebotHandler.RegisterBot()
	if err != nil {
		fmt.Printf("Error starting up bot: %s", err)
		return
	}

}
