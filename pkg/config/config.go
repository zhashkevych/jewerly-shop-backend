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

	if env := os.Getenv("HOST"); env == "stage" {
		viper.SetConfigName(env)
		return viper.MergeInConfig()
	}

	return nil
}