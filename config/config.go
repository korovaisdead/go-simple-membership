package config

import (
	"encoding/json"
	"os"
)

var (
	config *Configuration
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
		SaltLength int `json:"saltLength"`
		BcryptCost int `json:"bcryptCost"`
	} `json:"security"`
}

func GetConfig() (*Configuration, error) {
	if config != nil {
		return config, nil
	}

	file, err := os.Open("config.local.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
