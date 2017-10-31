package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	config *Configuration
)

type ConfigurationWeb struct {
	Port string `json:"port"`
}

type ConfigurationDb struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type ConfigurationSecurity struct {
	SaltLength int    `json:"saltLength"`
	BcryptCost int    `json:"bcryptCost"`
	Hmac       string `json:"hmac"`
}

type Configuration struct {
	Web      ConfigurationWeb      `json:"web"`
	Db       ConfigurationDb       `json:"db"`
	Security ConfigurationSecurity `json:"security"`
}

func GetConfig() (*Configuration, error) {
	if config != nil {
		return config, nil
	}
	return nil, errors.New("Please run the Build function before")
}

func Build(env string) (*Configuration, error) {
	filename := fmt.Sprintf("config.%v.json", env)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error importing the config file: ", filename)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding the configuration file")
		return nil, err
	}
	return config, nil
}

func BuildTestConfig() {
	if config == nil {
		config = &Configuration{
			Web: ConfigurationWeb{Port: ""},
			Db:  ConfigurationDb{Host: "localhost", Port: ":27018", Database: "test"},
		}
	}
}
