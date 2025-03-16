package main

import (
	"fmt"
	"os"

	"github.com/pjmessi/udp_chat/src/app"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	app.Run(username)
}
