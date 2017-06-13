package commons

import (
	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	Config = viper.New()
	Config.SetConfigName("app")
	Config.AddConfigPath("config/")
	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
}
