package conf

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	VbbApiUrl         string `json:"HAFAS_API_URL"`
	VbbAPIKey         string `json:"HAFAS_API_KEY"`
	MeteosourceApiKey string `json:"METEOSOURCE_API_KEY"`
	Coordinates       string `json:"COORDINATES"`
}

var Conf Config

func LoadConfig() {
	path := "conf/.conf.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("Config file not found:", path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Conf); err != nil {
		log.Fatal("Error decoding config file:", err)
	}

	if Conf.VbbApiUrl == "" || Conf.VbbAPIKey == "" || Conf.Coordinates == "" {
		log.Fatal("Missing required config keys")
	}
}
