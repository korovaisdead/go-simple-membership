package config

import (
	"errors"
	"os"
)

var (
	config *Configuration
)

type Configuration struct {
	WebPort             string
	DbUrl               string
	RedisUrl            string
	SecuritySecretWorld string
}

func GetConfig() (*Configuration, error) {
	if config != nil {
		return config, nil
	}
	return nil, errors.New("Please run the Build function before")
}

func Build() *Configuration {
	config = &Configuration{
		WebPort:             os.Getenv("WEB_PORT"),
		DbUrl:               os.Getenv("MONGO_URL"),
		RedisUrl:            os.Getenv("REDIS_URL"),
		SecuritySecretWorld: os.Getenv("SECURITY_SECRET_WORD"),
	}
	return config
}
