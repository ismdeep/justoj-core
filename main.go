package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := LoadConfig("config.example.json")
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(config.BaseUrl)

	fmt.Println(len(config.Languages))
}
