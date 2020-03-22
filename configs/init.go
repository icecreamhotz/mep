package configs

import (
	"github.com/spf13/viper"
)

func InitConfigs() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("configs")
	viper.ReadInConfig()

	viper.AddConfigPath("./configs")
	viper.SetConfigName("images")

	return viper.MergeInConfig()
}
