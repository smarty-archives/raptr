package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/smartystreets/raptr/cli"
	"github.com/smartystreets/raptr/config"
	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/messages"
	"github.com/smartystreets/raptr/storage"
	"github.com/smartystreets/raptr/tasks"
)

func main() {
	if command, configFile := cli.ReadMessage(); command == nil {
		log.Println("[ERROR] Unable to determine CLI command")
		os.Exit(1)
	} else if configuration, err := config.LoadConfiguration(configFile); err != nil {
		log.Println("[ERROR] Unable to load configuration:", err)
		os.Exit(1)
		// } else if err := configuration.Validate(command); err != nil {
		// 	log.Println("[ERROR] Invalid configuration:", err)
		// 	os.Exit(1)
	} else if err := executeCommand(configuration, command); err != nil {
		log.Println("[ERROR] Command Failed:", err)
		os.Exit(1)
	}
}
func executeCommand(configuration config.Configuration, message interface{}) error {
	switch message.(type) {
	case messages.UploadCommand:
		return executeUpload(configuration, message.(messages.UploadCommand))
	case messages.LinkCommand:
		return executeLink(configuration, message.(messages.LinkCommand))
	case messages.UnlinkCommand:
		return errors.New("Not implemented")
	case messages.CleanCommand:
		return errors.New("Not implemented")
	default:
		return errors.New("Not implemented")
	}
}
func executeUpload(configuration config.Configuration, command messages.UploadCommand) error {
	remote, found := configuration.Open(command.StorageName)
	if !found {
		return errors.New("Remote storage specified was not found in the configuration file.")
	}

	finder := manifest.NewLocalPackageFinder()
	task := tasks.NewUploadTask(remote.Storage)

	files, err := finder.Find(command.PackagePath)
	if err != nil {
		return err
	} else if len(files) == 0 {
		log.Println("[INFO] No files found; nothing to do.")
		return nil
	} else {
		// TODO: package name could be optional in one of the following scenarios:
		// 1. we have a DSC file
		// 2. there's only one package
		// otherwise, the name would be required.
		return task.Upload(command.Category, command.PackageName, files)
	}
}
func executeLink(configuration config.Configuration, command messages.LinkCommand) error {
	remote, found := configuration.Open(command.StorageName)
	if !found {
		return errors.New("Remote storage specified was not found in the configuration file.")
	}

	layout := remote.Layout
	task := tasks.NewLinkTask(remote.Storage, layout.Categories, layout.Architectures)

	for i := 0; i < maxConcurrencyAttempts; i++ {
		err := task.Link(command.Distribution, command.Category, command.PackageName, command.PackageVersion)
		if err == nil {
			return nil
		} else if err != storage.ConcurrencyError {
			return err
		}

		log.Printf("[WARN] Concurrency error, sleeping before trying again [%d/%d]\n", i+1, maxConcurrencyAttempts)
		time.Sleep(time.Second * 1)
	}

	return storage.ConcurrencyError
}

const maxConcurrencyAttempts = 5
