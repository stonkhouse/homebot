package config

type Configurations struct {
	Server   ServerConfigurations
	Database DBConfigurations
	Telebot  TelebotConfigurations
}

type ServerConfigurations struct {
	Port string
}

type DBConfigurations struct {
	Name     string
	Username string
	Password string
}

type TelebotConfigurations struct {
	Token string
}
