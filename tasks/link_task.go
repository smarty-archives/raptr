package tasks

import (
	"log"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type LinkTask struct {
	remote storage.Storage
	multi  *storage.MultiStorage
}

func NewLinkTask(remote storage.Storage) *LinkTask {
	return &LinkTask{remote: remote, multi: storage.NewMultiStorage(remote)}
}

// by this point, we've verified the category and target distribution(s)
func (this *LinkTask) Link(category, bundle, version string, distributions ...string) error {
	manifestPath := manifest.BuildPath(category, bundle, version)
	log.Println("[INFO] Downloading manifest from", manifestPath)
	response := this.remote.Get(storage.GetRequest{Path: manifestPath})
	manifestFile, err := manifest.ParseManifest(response.Contents, category, bundle, version)
	if err != nil {
		return err // unable to access or parse remote manifest, e.g. remote unavailable or permissions
	}

	state := NewIndexState(category, distributions, manifestFile.Architectures())
	gets := state.BuildGetRequests()
	if err := state.ReadGetResponses(this.multi.Get(gets...)); err != nil {
		return err // unable to access or parse remote Release|Sources|Packages file(s)
	}

	if err = state.Link(manifestFile); err != nil {
		return err
	}

	// TODO: GPG sign

	puts := state.BuildPutRequests()
	if err := state.ReadPutResponses(this.multi.Put(puts...)); err != nil {
		return err // concurrency, permissions, remote storage unavailable, etc.
	}

	return nil
}

func (this *LinkTask) buildGetRequests(category string, distributions, architectures []string) []storage.GetRequest {
	requests := []storage.GetRequest{}
	for _, distribution := range distributions {
		requests = append(requests, storage.GetRequest{Path: manifest.BuildReleaseFilePath(distribution)})
		for _, architecture := range architectures {
			path := ""
			if architecture == "source" {
				path = manifest.BuildSourcesFilePath(distribution, category)
			} else {
				path = manifest.BuildPackagesFilePath(distribution, category, architecture)
			}
			requests = append(requests, storage.GetRequest{Path: path})
		}
	}
	return requests
}
