package utils

import (
	"log"

	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/spf13/viper"
)

func LoadConfig(path string) *domain.Env {
	env := domain.Env{}
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
