package config

type Configurations struct {
	Server   ServerConfigurations
	Database DBConfigurations
}

type ServerConfigurations struct {
	Port string
}

type DBConfigurations struct {
	Name     string
	Username string
	Password string
}
