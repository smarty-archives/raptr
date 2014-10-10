package tasks

import (
	"crypto/md5"
	"errors"
	"io"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

// TODO: release files will need to know all architectures and categories
type IndexState struct {
	targetCategory string
	distributions  []string
	categories     []string
	architectures  []string
	items          []*IndexItem
}
type IndexItem struct {
	distribution       string
	targetArchitecture string
	previousMD5        []byte
	file               Serializable
}
type Serializable interface {
	Path() string
	Parse(io.Reader) error
	Bytes() []byte
}

func (this *IndexItem) IsReleaseFile() bool {
	return this.targetArchitecture == ""
}

func NewIndexState(targetCategory string, distributions, categories, architectures, targetArchitectures []string) *IndexState {
	this := &IndexState{
		targetCategory: targetCategory,
		distributions:  distributions,
		categories:     categories,
		architectures:  architectures,
		items:          []*IndexItem{},
	}

	for _, distribution := range distributions {
		this.items = append(this.items, &IndexItem{
			distribution: distribution,
			file:         manifest.NewReleaseFile(distribution, this.categories, this.architectures),
		})

		for _, targetArchitecture := range targetArchitectures {
			var file Serializable
			if targetArchitecture == "source" {
				file = manifest.NewSourcesFile(distribution, this.targetCategory)
			} else {
				file = manifest.NewPackagesFile(distribution, this.targetCategory, targetArchitecture)
			}

			this.items = append(this.items, &IndexItem{
				distribution:       distribution,
				targetArchitecture: targetArchitecture,
				file:               file,
			})
		}
	}

	return this
}

func (this *IndexState) BuildGetRequests() []storage.GetRequest {
	requests := []storage.GetRequest{}
	for _, item := range this.items {
		requests = append(requests, storage.GetRequest{Path: item.file.Path()})
	}
	return requests
}

func (this *IndexState) ReadGetResponses(responses []storage.GetResponse) error {
	if len(responses) != len(this.items) {
		return errors.New("Each request made should have exactly one response.")
	}

	for i, response := range responses {
		item := this.items[i]
		item.previousMD5 = response.MD5
		found := response.Error == nil
		if response.Error != nil && response.Error != storage.FileNotFoundError {
			return response.Error // only 404s are allowed here
		}

		if !found {
			continue
		} else if err := item.file.Parse(response.Contents); err != nil {
			return err
		}
	}

	return nil
}
func (this *IndexState) Link(file *manifest.ManifestFile) {
	releaseFiles := map[string]*manifest.ReleaseFile{}
	for _, item := range this.items {
		if item.IsReleaseFile() {
			releaseFiles[item.distribution] = item.file.(*manifest.ReleaseFile)
		} else {
			indexItem := item.file.(manifest.IndexFile)
			indexItem.Add(file)
			releaseFiles[item.distribution].Add(indexItem)
		}
	}
}
func (this *IndexState) GPGSign() error {
	return nil // TODO
}

func (this *IndexState) BuildPutRequests() []storage.PutRequest {
	requests := []storage.PutRequest{}
	for _, item := range this.items {
		payload := item.file.Bytes()
		md5sum := md5.Sum(payload)
		requests = append(requests, storage.PutRequest{
			Path:        item.file.Path(),
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
