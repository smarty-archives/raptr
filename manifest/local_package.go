package manifest

import "github.com/smartystreets/raptr/storage"

type LocalPackage interface {
	Name() string
	Version() string
	Architecture() string
	Metadata() Paragraph
	Files() []LocalPackageFile
}

type LocalPackageFile struct {
	Name      string
	Length    uint64
	Checksums Checksum
	Contents  storage.ReadSeekCloser
}
