package manifest

import (
	"errors"
	"io"
	"path"
)

type ManifestFile struct {
	category   string
	bundle     string
	version    string
	paragraphs []*Paragraph
	packages   map[string]struct{}
}

func NewManifestFile(category, bundle, version string) *ManifestFile {
	return &ManifestFile{
		category:   category,
		bundle:     bundle,
		version:    version,
		paragraphs: []*Paragraph{},
		packages:   map[string]struct{}{},
	}
}
func ParseManifest(reader io.Reader, category, bundle, version string) (*ManifestFile, error) {
	this := NewManifestFile(category, bundle, version)

	for {
		if paragraph, err := ReadParagraph(NewReader(reader)); err != nil {
			return nil, err
		} else if err == io.EOF {
			break
		} else if packageName, contains := paragraph.allKeys["Package"]; !contains {
			return nil, errors.New("Malformed manifest file, missing Package element.")
		} else if architecture, contains := paragraph.allKeys["Architecture"]; !contains {
			return nil, errors.New("Malformed manifest file, missing Architecture element.")
		} else if _, contains := paragraph.allKeys["Files"]; contains {
			this.packages[packageName.Value+"_source"] = struct{}{}
			this.paragraphs = append(this.paragraphs, paragraph)
		} else if _, contains := paragraph.allKeys["Filename"]; contains {
			this.packages[packageName.Value+"_"+architecture.Value] = struct{}{}
			this.paragraphs = append(this.paragraphs, paragraph)
		}
	}

	return this, nil
}
func BuildPath(category, bundle, version string) string {
	return path.Join("/pool/", category, bundle[0:1], bundle, version, "manifest") // FUTURE: gz?
}

func (this *ManifestFile) Path() string {
	return path.Dir(BuildPath(this.category, this.bundle, this.version))
}
func (this *ManifestFile) Add(pkg LocalPackage) (bool, error) {
	// TODO: only a single source package is allowed per manifest
	// check to see if it's already been added
	// watch out for the binary files with same package name and "_all" or "_any" prefix
	return false, nil
}
func (this *ManifestFile) Bytes() []byte {
	return []byte{} // TODO
}
