package config

import (
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
	return nil // TODO
}

func (this s3Config) buildStorage() (storage.Storage, error) {
	// FUTURE: from where else can/should we load security credentials?
	actual := storage.NewS3Storage(
		this.RegionName,
		this.BucketName,
		this.PathPrefix,
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"))

	if this.MaxRetries <= 0 {
		this.MaxRetries = defaultMaxRetries
	}

	inner := storage.Storage(actual)
	inner = storage.NewIntegrityStorage(inner)
	inner = storage.NewRetryStorage(inner, defaultMaxRetries)
	inner = storage.NewConcurrentStorage(inner)
	return inner, nil
}

const defaultMaxRetries = 3
