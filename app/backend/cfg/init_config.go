package cfg

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Server     ServerConfiguration
	CorsPolicy CorsPolicyConfiguration
	Database   DatabaseConfiguration
}

type ServerConfiguration struct {
	Verbose bool
	Port    int
}

type CorsPolicyConfiguration struct {
	Allow string
}

type DatabaseConfiguration struct {
	Username string
	Password string
	Address  string
	Port     string
}

var Config Configuration

func ReadConfigurationFile() error {
	viper.AddConfigPath("/opt/fci")
	viper.AddConfigPath("../")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return err
	}

	return nil
}
