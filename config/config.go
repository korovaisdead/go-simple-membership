package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	config *Configuration = nil
)

type Configuration struct {
	Web struct {
		Port string `json:"port"`
	} `json:"web"`
	Db struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	} `json:"db"`
	Security struct {
		SaltLength int    `json:"saltLength"`
		BcryptCost int    `json:"bcryptCost"`
		Hmac       string `json:"hmac"`
	} `json:"security"`
}

func BuildConfig(env string) (*Configuration, error) {
	if config != nil {
		fmt.Println("We have got the configturaion aleady")
		return config, nil
	}

	file, err := os.Open(fmt.Sprintf("config.%v.json", env))
	if err != nil {
		fmt.Println("Error importing the config file")
		return nil, err
	}
	defer file.Close()

	config = &Configuration{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		fmt.Println("Error decoding the configuration file")
		return nil, err
	}
	return config, nil
}

func Get() *Configuration {
	if &config == nil {
		panic("Forgot to init the configuration!")
	}
	return config
}
