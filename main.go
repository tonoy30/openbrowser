package main

import (
	"github.com/tonoy30/openbrowser/browser"
	"github.com/tonoy30/openbrowser/config"
	"log"
)

func main() {
	data := config.ParseJson()
	for _, val := range data.Urls {
		log.Printf("[Info] opening %s", val.Name)
		browser.OpenBrowser(val.URL)
	}

}
