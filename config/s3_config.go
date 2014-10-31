package config

import (
	"fmt"
	"os"

	"github.com/smartystreets/raptr/storage"
)

type s3Config struct {
	RegionName string `json:"region"`
	BucketName string `json:"bucket"`
	PathPrefix string `json:"prefix"`
	LayoutName string `json:"layout"`
	MaxRetries int    `json:"retries"`
	Timeout    int    `json:"timeout"`
}

func (this s3Config) validate() error {
	if this.BucketName == "" {
		return fmt.Errorf("The bucket name is missing.")
	} else {
		return nil
	}
}

func (this s3Config) buildStorage() (storage.Storage, error) {
	// FUTURE: from where else can/should we load security credentials?
	actual := storage.NewS3Storage(
		this.RegionName,
		this.BucketName,
		this.PathPrefix,
		getEnvironmentVariable("AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"),
		getEnvironmentVariable("AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"))

	if this.MaxRetries <= 0 {
		this.MaxRetries = defaultMaxRetries
	}

	inner := storage.Storage(actual)
	inner = storage.NewIntegrityStorage(inner)
	inner = storage.NewRetryStorage(inner, defaultMaxRetries)
	inner = storage.NewConcurrentStorage(inner)
	return inner, nil
}
func getEnvironmentVariable(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); len(value) > 0 {
			return value
		}
	}

	return ""
}

const defaultMaxRetries = 3
