package config

import (
	"github.com/sainak/status-checker/core/logger"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func GetConfig() *viper.Viper {
	if Config == nil {
		Config = viper.New()
		Config.AddConfigPath(".")
		Config.SetConfigName(".env")
		Config.SetConfigType("env")
		Config.AutomaticEnv()
		if err := Config.ReadInConfig(); err != nil {
			logger.Error(err)
		}
	}
	return Config
}

func GetDBurl() string {
	return Config.GetString("DB_URL")
}

func GetServerAddress() string {
	return ":" + Config.GetString("APP_PORT")
}
