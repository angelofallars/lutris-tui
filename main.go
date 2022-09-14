package main

import (
	"fmt"
	"log"
	view "lutris-tui/view"
	wrapper "lutris-tui/wrapper"
)

func main() {
	client, err := wrapper.NewLutrisClient()

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println("LOADING ...")
	games, err := client.FetchGames()

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	view.Start(client, games)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
