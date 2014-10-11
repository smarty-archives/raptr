package manifest

import "io"

type IndexFile interface {
	Path() string
	Add(*ManifestFile) bool

	Parse(io.Reader) error
	Bytes() []byte
}
