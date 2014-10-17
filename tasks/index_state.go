package tasks

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type IndexState struct {
	distribution           string
	availableCategories    []string
	availableArchitectures []string
	items                  []*IndexItem
}
type IndexItem struct {
	targetArchitecture string
	previousMD5        []byte
	file               IndexFile
}
type IndexFile interface {
	Path() string
	Parse(io.Reader) error
	Bytes() []byte
}

func NewIndexState(distribution string, availableCategories, availableArchitectures []string) *IndexState {
	return &IndexState{
		distribution:           distribution,
		availableCategories:    availableCategories,
		availableArchitectures: availableArchitectures,
		items: []*IndexItem{},
	}
}
func (this *IndexState) AddTarget(category string, architectures []string) {
	resolvedArchitectures := resolveArchitectures(this.availableArchitectures, architectures)
	log.Println("[INFO] Manifest contains packages with these CPU architectures:", resolvedArchitectures)

	releaseFile := manifest.NewReleaseFile(this.distribution, this.availableCategories, this.availableArchitectures)
	this.items = append(this.items, &IndexItem{file: releaseFile})

	for _, targetArchitecture := range resolvedArchitectures {
		this.items = append(this.items, &IndexItem{
			targetArchitecture: targetArchitecture,
			file:               buildIndexFile(this.distribution, category, targetArchitecture),
		})
	}
}
func resolveArchitectures(all, targets []string) []string {
	parsed := map[string]struct{}{}
	allowed := []string{}
	for _, target := range targets {
		if target == "all" || target == "any" {
			for _, item := range all {
				if item != "source" {
					parsed[item] = struct{}{}
				}
			}
		} else {
			parsed[target] = struct{}{}
		}
	}

	for key, _ := range parsed {
		allowed = append(allowed, key)
	}

	sort.Strings(allowed)
	return allowed
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
		log.Printf("[INFO] Downloading index file from %s.\n", item.file.Path())
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
			fmt.Println(item.file.Path())
			return err
		}
	}

	return nil
}
func (this *IndexState) Link(file *manifest.ManifestFile) bool {
	added := false
	var releaseFile *manifest.ReleaseFile
	for i, item := range this.items {
		if i == 0 {
			releaseFile = item.file.(*manifest.ReleaseFile)
		} else {
			indexItem := item.file.(manifest.IndexFile)
			if indexItem.Add(file) {
				added = true
				releaseFile.Add(indexItem)
			}
		}
	}
	added = true // for right now, always publish even if nothing's changed
	return added
}
func (this *IndexState) GPGSign() error {
	releaseFile := this.items[0].file.Bytes()
	if signature, err := SignDistributionIndex(this.distribution, releaseFile); err != nil {
		return err
	} else {
		this.items = append(this.items, &IndexItem{file: signature})
		return nil
	}
}

func (this *IndexState) BuildPutRequests() []storage.PutRequest {
	requests := []storage.PutRequest{}
	for _, item := range this.items {
		payload := item.file.Bytes()
		md5sum := md5.Sum(payload)
		log.Printf("[INFO] Uploading index to %s.\n", item.file.Path())
		requests = append(requests, storage.PutRequest{
			Path:        item.file.Path(),
			Length:      uint64(len(payload)),
			Contents:    storage.NewReader(payload),
			MD5:         md5sum[:],
			ExpectedMD5: item.previousMD5, // make sure nothing has changed
			Concurrency: storage.CheckBeforePut | storage.CheckAfterPut,
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
