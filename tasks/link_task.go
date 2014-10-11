package tasks

import (
	"fmt"
	"log"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type LinkTask struct {
	remote        storage.Storage
	multi         *storage.MultiStorage
	categories    []string
	architectures []string
}

func NewLinkTask(remote storage.Storage, categories, architectures []string) *LinkTask {
	return &LinkTask{
		remote:        remote,
		multi:         storage.NewMultiStorage(remote),
		categories:    categories,
		architectures: architectures,
	}
}

// by this point, we've verified the category and target distribution(s)
func (this *LinkTask) Link(distribution, category, bundle, version string) error {
	manifestPath := manifest.BuildPath(category, bundle, version)
	log.Println("[INFO] Downloading manifest from", manifestPath)
	response := this.remote.Get(storage.GetRequest{Path: manifestPath})
	if response.Error == storage.FileNotFoundError {
		return fmt.Errorf("No manifest file exists for bundle [%s_%s] in [%s].", bundle, version, category)
	}

	manifestFile, err := manifest.ParseManifest(response.Contents, category, bundle, version)
	if err != nil {
		return err // unable to access or parse remote manifest, e.g. remote unavailable or permissions
	}

	log.Println("[INFO] Manifest parsed.")
	state := NewIndexState(distribution, this.categories, this.architectures)
	state.AddTarget(category, manifestFile.Architectures())
	gets := state.BuildGetRequests()
	if err := state.ReadGetResponses(this.multi.Get(gets...)); err != nil {
		log.Println("[ERROR] Failed to parse one or more index files.")
		return err // unable to access or parse remote Release|Sources|Packages file(s)
	}

	log.Println("[INFO] Linking manifest to downloaded indexes.")
	if !state.Link(manifestFile) {
		log.Println("[INFO] Manifest already found in index files; no changes performed.")
		return nil
	}

	if err = state.GPGSign(); err != nil {
		log.Println("[ERROR] Unable to sign the Release file.")
		return err
	}

	puts := state.BuildPutRequests()
	if err := state.ReadPutResponses(this.multi.Put(puts...)); err != nil {
		return err // concurrency, permissions, remote storage unavailable, etc.
	}

	return nil
}
