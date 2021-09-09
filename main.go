package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/tucnak/telebot.v2"
	"stonkhouse/stonkbot/bot"
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
	homebot, err := telebot.NewBot(telebot.Settings{
		Token:  config.Telebot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	bot.RegisterBot(homebot)
	if err != nil {
		fmt.Printf("Error starting up bot: %s", err)
		return
	}

}
