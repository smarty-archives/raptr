package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/raptr/config"
	"github.com/smartystreets/raptr/storage"
)

func main() {
	if configuration, err := config.LoadConfiguration("raptr.conf"); err != nil {
		log.Println("[CONFIG ERROR]:", err)
		os.Exit(1)
	} else if repo, found := configuration.Open("repo-1"); !found {
		log.Println("[CONFIG ERROR]: Repo named 'repo-1' not found.")
		os.Exit(1)
	} else {
		expected := []byte{}
		expected, _ = hex.DecodeString("b4ae6236dedc23dc45396c33e6550fb0")
		response := repo.Storage.Head(storage.HeadRequest{Path: "/pubkey.asc", ExpectedMD5: expected})
		fmt.Println("Expected MD5:", hex.EncodeToString(expected))
		fmt.Println("Actual MD5:", hex.EncodeToString(response.MD5))
		fmt.Println("Length:", response.Length)
		fmt.Println("Error:", response.Error)
	}
}
