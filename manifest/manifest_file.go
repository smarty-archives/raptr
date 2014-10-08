package manifest

import (
	"io"
	"path"
)

type ManifestFile struct {
	category string
	bundle   string
	version  string
}

func NewManifestFile(category, bundle, version string) *ManifestFile {
	return &ManifestFile{category: category, bundle: bundle, version: version}
}
func ParseManifest(reader io.Reader) (*ManifestFile, error) {
	return nil, nil
}
func BuildPath(category, bundle, version string) string {
	return path.Join("/pool/", category, bundle[0:1], bundle, version, "manifest") // FUTURE: gz?
}

func (this *ManifestFile) Path() string {
	return path.Base(BuildPath(this.category, this.bundle, this.version))
}
func (this *ManifestFile) Add(pkg LocalPackage) (bool, error) {
	return false, nil
}
func (this *ManifestFile) Bytes() []byte {
	return []byte{} // TODO
}
