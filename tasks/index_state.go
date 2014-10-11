package tasks

import (
	"crypto/md5"
	"errors"
	"io"
	"log"
	"sort"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type IndexState struct {
	targetCategory   string
	distributions    []string
	allCategories    []string
	allArchitectures []string
	items            []*IndexItem
}
type IndexItem struct {
	distribution       string
	targetArchitecture string
	previousMD5        []byte
	file               IndexFile
}
type IndexFile interface {
	Path() string
	Parse(io.Reader) error
	Bytes() []byte
}

func NewIndexState(targetCategory string, distributions, allCategories, allArchitectures, targetArchitectures []string) *IndexState {
	this := &IndexState{
		targetCategory:   targetCategory,
		distributions:    distributions,
		allCategories:    allCategories,
		allArchitectures: allArchitectures,
		items:            []*IndexItem{},
	}

	targetArchitectures = findTargetArchitectures(allArchitectures, targetArchitectures)
	log.Println("[INFO] Manifest contains packages with these architectures:", targetArchitectures)

	for _, distribution := range distributions {
		this.items = append(this.items, &IndexItem{
			distribution: distribution,
			file:         manifest.NewReleaseFile(distribution, this.allCategories, this.allArchitectures),
		})

		for _, targetArchitecture := range targetArchitectures {
			this.items = append(this.items, &IndexItem{
				distribution:       distribution,
				targetArchitecture: targetArchitecture,
				file:               buildIndexFile(distribution, this.targetCategory, targetArchitecture),
			})
		}
	}

	return this
}
func findTargetArchitectures(all, targets []string) []string {
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
		log.Printf("[INFO] Downloading file from %s.\n", item.file.Path())
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
		if item.targetArchitecture == "" {
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
		log.Printf("[INFO] Uploading index to %s.\n", item.file.Path())
		requests = append(requests, storage.PutRequest{
			Path:        item.file.Path(),
			Length:      uint64(len(payload)),
			Contents:    storage.NewReader(payload),
			MD5:         md5sum[:],
			ExpectedMD5: item.previousMD5, // make sure nothing has changed
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
