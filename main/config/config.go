package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once
)

// Init инициализирует Viper и загружает конфигурацию
func Init() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
	})
}

func GetString(key string) string {
	Init()
	return viper.GetString(key)
}

func GetInt(key string) int {
	Init()
	return viper.GetInt(key)
}
