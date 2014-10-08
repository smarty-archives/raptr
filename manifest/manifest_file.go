package manifest

import (
	"io"
	"path"
)

type ManifestFile struct{}

func NewManifestFile() *ManifestFile {
	return &ManifestFile{}
}
func ParseManifest(reader io.Reader) (*ManifestFile, error) {
	return nil, nil
}
func BuildPath(category, bundle, version string) string {
	return path.Join("/", category, bundle, version, "manifest") // TODO: gz?
}

func (this *ManifestFile) Add(filename string) bool {
	return false // TODO
}
func (this *ManifestFile) Bytes() []byte {
	return []byte{} // TODO
}
