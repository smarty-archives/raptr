package manifest

import (
	"errors"
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

func (this *SourcesFile) Add(manifest *ManifestFile) bool {
	added := false
	for _, paragraph := range manifest.architectures["source"] {
		name, version := paragraph.Name(), paragraph.Version()
		id := name + "_" + version
		if id == "_" {
			continue // bad paragraph
		} else if _, contains := this.packages[id]; contains {
			continue // already exists
		} else {
			added = true
			this.cachedBytes = nil
			this.packages[id] = struct{}{}
			this.paragraphs = append(this.paragraphs, paragraph)
		}
	}

	return added
}

func (this *SourcesFile) Parse(reader io.Reader) error {
	this.cachedBytes = nil
	this.paragraphs = []*Paragraph{}
	this.packages = map[string]struct{}{}

	paragraphReader := NewReader(reader)

	for {
		if paragraph, err := ReadParagraph(paragraphReader); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if name, contains := paragraph.allKeys["Package"]; !contains {
			return errors.New("Malformed manifest file, missing Package element.")
		} else if version, contains := paragraph.allKeys["Version"]; !contains {
			return errors.New("Malformed manifest file, missing Version element.")
		} else {
			this.packages[name.Value+"_"+version.Value] = struct{}{}
			this.paragraphs = append(this.paragraphs, paragraph)
		}
	}

	return nil
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
