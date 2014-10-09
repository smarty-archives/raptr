package main

import (
	"fmt"

	"github.com/smartystreets/raptr/manifest"
)

func main() {
	directory := "/Users/jonathan/Downloads/tmp"
	finder := manifest.NewLocalPackageFinder()

	manifestFile := manifest.NewManifestFile("public", "nginx", "1.7.4-1")
	if files, err := finder.Find(directory); err != nil {
		fmt.Println(err)
	} else {
		for _, file := range files {
			if success, err := manifestFile.Add(file); err != nil {
				fmt.Println(err)
			} else if !success {
				fmt.Println("failed")
			}
		}
	}

	fmt.Print(string(manifestFile.Bytes()))
}
