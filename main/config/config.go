package config

import (
	"sync"

	"github.com/spf13/viper"

)

var (
	once sync.Once
)

// Init инициализирует Viper и загружает конфигурацию
func init() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			panic(err.Error())
		}
	})
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}
