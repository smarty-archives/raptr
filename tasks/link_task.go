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

	requests := this.buildGetRequests(category, distributions, manifestFile.Architectures())
	responses := this.multi.Get(requests...)
	for _, response := range responses {
		if response.Error != nil && response.Error != storage.FileNotFoundError {
			return response.Error
		}

		// three kinds of files--Release, Packages, Sources; parse each one as the appropriate type
		// create an "index" out of it, e.g. a ReleasesFile and SourcesFile implement the IndexFile interface
		// on each one, call: index.Add(manifestFile) // which adds just the bits that it needs
		// on the root one call release.Add(index) // which computes the various hashes of the bytes
	}

	// sign Releases file (GPG) (do this last)
	// upload all files (Packages|Sources|Release)--pass any concurrency errors up the chain
	//    to the controlling code (which should re-run this task)
	return nil
}
func (this *LinkTask) buildGetRequests(category string, distributions, architectures []string) []storage.GetRequest {
	requests := []storage.GetRequest{storage.GetRequest{Path: manifest.BuildReleaseFilePath()}}
	for _, distribution := range distributions {
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
