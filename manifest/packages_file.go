package manifest

import (
	"errors"
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
	return path.Join("/dists/", distribution, category, "binary-"+architecture, "Packages")
}

func (this *PackagesFile) Add(manifest *ManifestFile) bool {
	added := false
	for _, architecture := range []string{this.architecture, "any", "all"} {
		for _, paragraph := range manifest.architectures[architecture] {
			name, version := paragraph.PackageName(), paragraph.Version()
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
	}

	return added
}

func (this *PackagesFile) Parse(reader io.Reader) error {
	this.cachedBytes = nil
	this.paragraphs = []*Paragraph{}
	this.packages = map[string]struct{}{}

	// gzipReader, err := gzip.NewReader(reader)
	// if err != nil {
	// 	return err
	// }
	// paragraphReader := NewReader(gzipReader)
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

func (this *PackagesFile) Bytes() []byte {
	if this.cachedBytes != nil {
		return this.cachedBytes
	}

	this.cachedBytes = compressAndSerializeParagraphs(this.paragraphs)
	return this.cachedBytes
}

func (this *PackagesFile) Path() string {
	return this.path
}
