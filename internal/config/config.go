package config

import (
	"os"
	"fmt"
	"encoding/json"
	"log"
)

type Config struct {
	DBUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c * Config) SetUser(user string) {
	homeDir, err := os.UserHomeDir()

	c.CurrentUserName = user
	configFile, err := os.OpenFile(fmt.Sprintf("%s/.gatorconfig.json", homeDir), os.O_RDWR|os.O_CREATE, 0644)
	defer configFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	configEncoder := json.NewEncoder(configFile)
	if err := configEncoder.Encode(&c); err != nil {
		log.Fatal(err)
	}
}

func Read() Config {
	homeDir, err := os.UserHomeDir()
	
	if err != nil {
		log.Fatal(err)
	}

	var c Config
	configFile, err := os.Open(fmt.Sprintf("%s/.gatorconfig.json", homeDir))
	defer configFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	configDecoder := json.NewDecoder(configFile)
	if err := configDecoder.Decode(&c); err != nil {
		log.Fatal(err)
	}
	return Config { DBUrl: c.DBUrl, CurrentUserName: c.CurrentUserName }
	
}
