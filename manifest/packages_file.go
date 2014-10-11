package manifest

import (
	"io"
	"path"
)

// Represents a collection of compiled debian binaries that share the same CPU architecture
// NOTE: we will read and write these files
type PackagesFile struct {
	cachedBytes  []byte
	path         string
	architecture string
	paragraphs   []*Paragraph
	packages     map[string]struct{}
}

func NewPackagesFile(distribution, category, architecture string) *PackagesFile {
	return &PackagesFile{
		path:         BuildPackagesFilePath(distribution, category, architecture),
		architecture: architecture,
		paragraphs:   []*Paragraph{},
		packages:     map[string]struct{}{},
	}
}
func BuildPackagesFilePath(distribution, category, architecture string) string {
	return path.Join("/dists/", distribution, category, "binary-"+architecture, "Packages.gz")
}

func (this *PackagesFile) Add(manifest *ManifestFile) {
	this.cachedBytes = nil

	for _, architecture := range []string{this.architecture, "any", "all"} {
		for _, paragraph := range manifest.architectures[architecture] {
			name, version := paragraph.Name(), paragraph.Version()
			id := name + "_" + version
			if id == "_" {
				continue // bad paragraph
			} else if _, contains := this.packages[id]; contains {
				continue // already exists
			} else {
				this.packages[id] = struct{}{}
				this.paragraphs = append(this.paragraphs, paragraph)
			}
		}
	}
}

func (this *PackagesFile) Parse(reader io.Reader) error {
	return nil // TODO
}

func (this *PackagesFile) Bytes() []byte {
	if this.cachedBytes != nil {
		return this.cachedBytes
	}

	this.cachedBytes = serializeParagraphs(this.paragraphs)
	return this.cachedBytes
}

func (this *PackagesFile) Path() string {
	return this.path
}
