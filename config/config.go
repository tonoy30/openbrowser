package config

import (
	"encoding/json"
	"log"
	"os"
)

type Urls struct {
	Urls []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"urls"`
}

func ParseJson() *Urls {
	var urls Urls
	file, err := os.Open("./config/urls.json")
	if err != nil {
		log.Fatalf("%s error in reading json file")
	}
	defer file.Close()
	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(&urls)
	return &urls
}
