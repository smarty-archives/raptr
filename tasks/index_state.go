package tasks

import (
	"crypto/md5"
	"errors"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type IndexState struct {
	category      string
	distributions []string
	architectures []string
	items         []IndexItem
	releaseFiles  map[string]*manifest.ReleaseFile
	indexFiles    map[string]manifest.IndexFile
}
type IndexItem struct {
	path         string
	distribution string
	architecture string
	previousMD5  []byte
	file         Serializable
}
type Serializable interface {
	Bytes() []byte
}

func NewIndexState(category string, distributions, architectures []string) *IndexState {
	this := &IndexState{
		category:      category,
		distributions: distributions,
		architectures: architectures,
		items:         []IndexItem{},
		releaseFiles:  map[string]*manifest.ReleaseFile{},
		indexFiles:    map[string]manifest.IndexFile{},
	}

	for _, distribution := range distributions {
		this.items = append(this.items, IndexItem{
			path:         manifest.BuildReleaseFilePath(distribution),
			distribution: distribution,
		})

		for _, architecture := range architectures {
			path := ""
			if architecture == "source" {
				path = manifest.BuildSourcesFilePath(distribution, this.category)
			} else {
				path = manifest.BuildPackagesFilePath(distribution, this.category, architecture)
			}

			this.items = append(this.items, IndexItem{
				path:         path,
				distribution: distribution,
				architecture: architecture,
			})
		}
	}

	return this
}
func buildIndexFile(distribution, category, architecture string) manifest.IndexFile {
	if architecture == "source" {
		return manifest.NewSourcesFile(distribution, category)
	} else {
		return manifest.NewPackagesFile(distribution, category, architecture)
	}
}

func (this *IndexState) BuildGetRequests() []storage.GetRequest {
	requests := []storage.GetRequest{}
	for _, item := range this.items {
		requests = append(requests, storage.GetRequest{Path: item.path})
	}
	return requests
}

func (this *IndexState) ReadGetResponses(responses []storage.GetResponse) error {
	if len(responses) != len(this.items) {
		return errors.New("Each request made should have exactly one response.")
	}

	// releaseFile := manifest.NewReleaseFile(distribution, categories, architectures)
	// this.releaseFiles = append(this.releasesFile, releaseFile)

	// indexFile := buildIndexFile(distribution, category, architecture)
	// this.indexFiles = append(this.indexFiles, indexFile)

	for i, response := range responses {
		if response.Error != nil && response.Error != storage.FileNotFoundError {
			return response.Error // only 404s are allowed here
		}

		this.items[i].previousMD5 = response.MD5
	}

	return nil
}
func (this *IndexState) Link(file *manifest.ManifestFile) error {
	return nil
}

func (this *IndexState) BuildPutRequests() []storage.PutRequest {
	requests := []storage.PutRequest{}
	for _, item := range this.items {
		payload := item.file.Bytes()
		md5sum := md5.Sum(payload)
		requests = append(requests, storage.PutRequest{
			Path:        item.path,
			MD5:         md5sum[:],
			ExpectedMD5: item.previousMD5, // make sure nothing has changed
			Length:      uint64(len(payload)),
		})
	}
	return requests
}
func (this *IndexState) ReadPutResponses(responses []storage.PutResponse) error {
	if len(responses) != len(this.items) {
		return errors.New("Each request made should have exactly one response.")
	}

	for _, response := range responses {
		if response.Error != nil {
			return response.Error
		}
	}

	return nil
}
