package main

import (
	"fmt"
	"github.com/ismdeep/justoj-core/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig("config.example.json")
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(config.BaseUrl)

	fmt.Println(len(config.Languages))
}
