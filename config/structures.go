package config

import "github.com/smartystreets/raptr/storage"

type ConfigFormat struct {
	Layouts map[string]LayoutInfo
	S3      map[string]S3Info `json:"s3"`
}

// TODO: function to validate contents
type LayoutInfo struct {
	LayoutKey     string   `json:-`
	Distributions []string `json:"distributions"`
	Categories    []string `json:"categories"`
	Architectures []string `json:"architectures"`
}

// TODO: function to validate contents
type S3Info struct {
	StorageKey string `json:-`
	RegionName string `json:"region"`
	BucketName string `json:"bucket"`
	PathPrefix string `json:"prefix"`
	LayoutName string `json:"layout"`
	MaxRetries int    `json:"retries"`
	Timeout    int    `json:"timeout"`
}

type RepositoryInfo struct {
	StorageKey string
	Storage    storage.Storage
	Layout     LayoutInfo
}

func (this S3Info) Validate() error {
	return nil
}
func (this LayoutInfo) Validate() error {
	return nil
}
