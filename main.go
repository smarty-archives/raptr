package main

import (
	"fmt"

	"github.com/smartystreets/raptr/config"
	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/tasks"
)

func main() {
	configuration, _ := config.LoadConfiguration("raptr.conf")
	repo, _ := configuration.Open("repo-1")

	directory := "/home/vagrant/tmp"
	finder := manifest.NewLocalPackageFinder()
	task := tasks.NewUploadTask(repo.Storage)

	if files, err := finder.Find(directory); err != nil {
		fmt.Println(err)
	} else if err := task.Upload("public", "nginx", "1.1.19-1ubuntu0.6", files); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success!")
	}
}
