package manifest

import (
	"io"
	"path"
)

// Represents a set of debian source code archives for a given software category
// NOTE: we will read and write these files
type SourcesFile struct {
	cachedBytes []byte
	path        string
	paragraphs  []*Paragraph
	packages    map[string]struct{}
}

func NewSourcesFile(distribution, category string) *SourcesFile {
	return &SourcesFile{
		path:       BuildSourcesFilePath(distribution, category),
		paragraphs: []*Paragraph{},
		packages:   map[string]struct{}{},
	}
}
func BuildSourcesFilePath(distribution, category string) string {
	return path.Join("/dists/", distribution, category, "source/Sources")
}

func (this *SourcesFile) Add(manifest *ManifestFile) {
	this.cachedBytes = nil
	for _, paragraph := range manifest.architectures["source"] {
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

func (this *SourcesFile) Parse(reader io.Reader) error {
	return nil // TODO
}

func (this *SourcesFile) Bytes() []byte {
	if this.cachedBytes != nil {
		return this.cachedBytes
	}

	this.cachedBytes = serializeParagraphs(this.paragraphs)
	return this.cachedBytes
}

func (this *SourcesFile) Path() string {
	return this.path
}
