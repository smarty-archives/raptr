package manifest

import "github.com/smartystreets/raptr/storage"

type LocalPackage interface {
	Name() string
	Version() string
	Architecture() string
	Files() []LocalPackageFile
}

type LocalPackageFile struct {
	Name     string
	MD5      []byte
	Contents storage.ReadSeekCloser
}
