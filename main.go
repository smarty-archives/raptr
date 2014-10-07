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
	fmt.Println("Configuration:", configuration)

	switch command.(type) {
	case messages.UploadCommand:
		fmt.Println("Upload...", command)
	case messages.LinkCommand:
		fmt.Println("Link...", command)
	case messages.UnlinkCommand:
		fmt.Println("Unlink...", command)
	case messages.CleanCommand:
		fmt.Println("Clean...", command)
	default:
		fmt.Printf("Unknown command\n")
	}

}
