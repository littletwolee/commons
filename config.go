package commons

import (
	"github.com/spf13/viper"
)

var config *viper.Viper

func GetConfig() *viper.Viper {
	if config == nil {
		config = viper.New()
		config.SetConfigName("app")
		config.AddConfigPath("config/")
		if err := config.ReadInConfig(); err != nil {
			panic(err)
		}
	}
	return config
}
