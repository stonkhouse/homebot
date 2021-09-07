package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	c "stonkhouse/stonkbot/config"
	"stonkhouse/stonkbot/healthcheck"
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
	mainRouter := gin.Default()
	mainRouter.GET("", healthcheck.GetHealth)
	err := mainRouter.Run(":" + config.Server.Port)
	if err != nil {
		return
	}
}
