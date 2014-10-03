package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/smartystreets/raptr/storage"
)

type Configuration struct {
	repos map[string]RepositoryInfo
}

func LoadConfiguration(fullPath string) (Configuration, error) {
	fullPath = filepath.Clean(fullPath)

	// TODO
	// fullPath (if provided)
	// working directory/.ratpr
	// home directory/.raptr
	// /etc/raptr.conf
	return readFile(fullPath)
}
func (this Configuration) Repository(name string) (RepositoryInfo, bool) {
	item, found := this.repos[name]
	return item, found
}

func readFile(fullPath string) (Configuration, error) {
	deserialized := ConfigFormat{}

	if handle, err := os.Open(fullPath); err != nil {
		return Configuration{}, err // file doesn't exist or access problems
	} else if payload, err := ioutil.ReadAll(handle); err != nil {
		handle.Close()
		return Configuration{}, err // couldn't read file
	} else if err := json.Unmarshal(payload, &deserialized); err != nil {
		handle.Close()
		return Configuration{}, err // malformed JSON
	} else {
		handle.Close()
		return newConfiguration(deserialized)
	}
}
func newConfiguration(format ConfigFormat) (Configuration, error) {
	repos := map[string]RepositoryInfo{}
	layouts := map[string]LayoutInfo{}

	for key, item := range format.Layouts {
		item.LayoutKey = key
		if err := item.Validate(); err != nil {
			return Configuration{}, fmt.Errorf("Layout '%s' has missing or corrupt values.", key)
		} else {
			layouts[key] = item
		}
	}

	for key, item := range format.S3 {
		item.StorageKey = key
		if layout, found := layouts[item.LayoutName]; !found {
			return Configuration{}, fmt.Errorf("S3 store '%s' references not-existent layout '%s'.", key, item.LayoutName)
		} else if err := item.Validate(); err != nil {
			return Configuration{}, fmt.Errorf("S3 store '%s' has missing or corrupt values.", key)
		} else if store, err := newS3Storage(item); err != nil {
			return Configuration{}, fmt.Errorf("S3 store '%s' cannot be initialized.", key)
		} else {
			repos[key] = RepositoryInfo{StorageKey: key, Storage: store, Layout: layout}
		}
	}

	return Configuration{repos: repos}, nil
}
func newS3Storage(info S3Info) (storage.Storage, error) {
	// TODO: from where else can/should we load security credentials?
	actual := storage.NewS3Storage(
		info.RegionName,
		info.BucketName,
		info.PathPrefix,
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"))

	if info.MaxRetries <= 0 {
		info.MaxRetries = defaultMaxRetries
	}

	inner := storage.Storage(actual)
	inner = storage.NewIntegrityStorage(inner)
	inner = storage.NewRetryStorage(inner, defaultMaxRetries)
	inner = storage.NewConcurrentStorage(inner)
	return inner, nil
}

const defaultMaxRetries = 3
