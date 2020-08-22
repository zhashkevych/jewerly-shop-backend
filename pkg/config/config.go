package config

import (
	"github.com/spf13/viper"
	"os"
)

func Init() error {
	viper.AddConfigPath("./config")

	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	env := os.Getenv("HOST")

	switch env {
	case "stage", "prod":
		viper.SetConfigName(env)
		return viper.MergeInConfig()
	default:
		return nil
	}
}
