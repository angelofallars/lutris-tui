package main

import (
	"fmt"
	"log"
	view "lutris-tui/view"
	wrapper "lutris-tui/wrapper"
)

func main() {
	lutris, err := wrapper.NewWrapper()

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println("LOADING ...")
	games, err := lutris.FetchGames()

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	view.Start(lutris, games)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
