package main

import (
	"fmt"

	"github.com/smartystreets/raptr/cli"
	"github.com/smartystreets/raptr/config"
	"github.com/smartystreets/raptr/messages"
)

func main() {
	command, configFile := cli.ReadMessage()
	configuration, _ := config.LoadConfiguration(configFile)

	switch command.(type) {
	case messages.UploadCommand:
		fmt.Println("Upload...", configuration)
	case messages.LinkCommand:
		fmt.Println("Link...", configuration)
	case messages.UnlinkCommand:
		fmt.Println("Unlink...", configuration)
	case messages.CleanCommand:
		fmt.Println("Clean...", configuration)
	default:
		fmt.Printf("Unknown command\n")
	}

}
