package config

type Configurations struct {
	Server   ServerConfigurations
	Firebase DBConfigurations
	Telebot  TelebotConfigurations
}

type ServerConfigurations struct {
	Port string
}

type DBConfigurations struct {
	ConfigString string
}

type TelebotConfigurations struct {
	Token string
}
