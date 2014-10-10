package manifest

import "io"

type IndexFile interface {
	Path() string
	Add(*ManifestFile)

	Parse(io.Reader) error
	Bytes() []byte
}
