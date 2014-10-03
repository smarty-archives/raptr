package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/smartystreets/raptr/config"
	"github.com/smartystreets/raptr/storage"
)

func main() {
	usage := `
Usage:
    raptr upload
                  --name=PACKAGE
                  --path=PATH
                  --storage=STORE
                  --category=CATEGORY
                  [--distribution=DISTRIBUTION]
    raptr link
                  --name=PACKAGE
                  --version=VERSION
                  --storage=named-STORE
                  --category=CATEGORY
                  --source-distribution=DISTRIBUTION
                  --target-distribution=DISTRIBUTION
    ratpr unlink
                  --name=PACKAGE
                  --version=VERSION
                  --storage=named-STORE
                  --category=CATEGORY
                  --distribution=DISTRIBUTION
    raptr clean --storage=STORE
    raptr -h | -v

Examples:
    raptr upload --name=liveaddress-logging --path=/package/files --storage=s3 --category=operations [--distribution=staging]
    raptr link --name=liveaddress-logging --version=1.0.7 --storage=named-storage --category=liveaddress --source-distribution=staging --target-distribution=production
    ratpr unlink --name=liveaddress-logging --version=1.0.7 --storage=named-storage --category=liveaddress --distribution=production
    raptr clean --storage=named-storage


Options:
    -h  Shows this screen
    -v  Shows the version

    --name=PACKAGE                      Name of the package
    --path=PATH                         Path to package files
    --storage=STORE                     Storage mechanism to use (eg: s3, filesystem, etc)
    --category=CATEGORY                 Category the package falls under
    --distribution=DISTRIBUTION         Distribution to link to (optional for upload, if exists, links package to distribution)
    --source-distribution=DISTRIBUTION	Source distribution
    --target-distribution=DISTRIBUTION  Target distribution
`
	//docopt.Parse(doc, argv, help, version, optionsFirst, ...)
	arguments, err := docopt.Parse(usage, nil, true, "raptr (Remote Apt Repository) version 1.2.3", false)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("arguments:", arguments)
	for key, val := range arguments {
		fmt.Printf("%v - %v\n", key, val)
	}
	if arguments["upload"] == true {
		fmt.Println("We have an upload request:")
		fmt.Println("name:", arguments["--name"])
		fmt.Println("path:", arguments["--path"])
		fmt.Println("storage:", arguments["--storage"])
		fmt.Println("category:", arguments["--category"])
		fmt.Println("")
	}

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
