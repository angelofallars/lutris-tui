package main

import (
	"fmt"
	"log"
	"lutris-tui/lutris"
	view "lutris-tui/view"
)

func main() {
	client, err := lutris.NewLutrisClient()

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
