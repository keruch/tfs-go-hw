package config

import (
	"net/url"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/spf13/viper"
)

func SetupConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./trading_robot")
	viper.AddConfigPath("./trading_robot/config")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetPrivateKey() string {
	return viper.GetString("API.private_key")
}

func GetPublicKey() string {
	return viper.GetString("API.public_key")
}

func GetPeriod() domain.CandlePeriod {
	return domain.CandlePeriod(viper.GetString("pair.period"))
}

func GetDatabaseURL() string {
	u := url.URL{
		Host:   viper.GetString("database.address") + viper.GetString("database.port"),
		User:   url.UserPassword(viper.GetString("database.username"), viper.GetString("database.password")),
		Scheme: viper.GetString("database.scheme"),
		Path:   viper.GetString("database.name"),
	}
	return u.String()
}

func GetTelegramBotToken() string {
	return viper.GetString("API.tg_bot_token")
}

func GetServerAddress() string {
	return viper.GetString("server.address")
}
