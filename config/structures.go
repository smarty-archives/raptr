package config

import "github.com/smartystreets/raptr/storage"

type configFile struct {
	Layouts map[string]RepositoryLayout
	S3      map[string]s3Config `json:"s3"`
}
type RepositoryConfig struct {
	Storage storage.Storage
	Layout  RepositoryLayout
}
