package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/smartystreets/raptr/storage"
)

type Configuration struct {
	repos map[string]RepositoryInfo
}

func LoadConfiguration(fullPath string) (Configuration, error) {
	// TODO
	// fullPath (if provided)
	// working directory/.ratpr
	// home directory/.raptr
	// /etc/raptr.conf
	return readFile(fullPath)
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
		return newConfiguration(deserialized), nil
	}
}

func newConfiguration(format ConfigFormat) Configuration {
	repos := map[string]RepositoryInfo{}
	layouts := map[string]LayoutInfo{}

	for key, item := range format.Layouts {
		item.LayoutKey = key
		layouts[key] = item
	}

	for key, item := range format.S3 {
		item.StorageKey = key

		if layout, found := layouts[item.LayoutName]; !found {
			log.Printf("[CONFIG] S3 store '%s' references not-existent layout '%s'.\n", key, item.LayoutName)
			os.Exit(1)
		} else if store, err := newS3Storage(item); err != nil {
			log.Printf("[CONFIG] S3 store '%s' cannot be initialized.\n", key)
			os.Exit(1)
		} else {
			repos[key] = RepositoryInfo{
				StorageKey: key,
				Storage:    store,
				Layout:     layout,
			}
		}

	}

	return Configuration{}
}
func newS3Storage(info S3Info) (storage.Storage, error) {
	// TODO: missing/empty values in the configuration file, e.g. bucket name
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
