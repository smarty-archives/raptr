package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
)

type ManifestFile struct {
	category   string
	bundle     string
	version    string
	hasDSC     bool
	paragraphs []*Paragraph
	packages   map[string]struct{}
}

func NewManifestFile(category, bundle, version string) *ManifestFile {
	return &ManifestFile{
		category:   category,
		bundle:     bundle,
		version:    version,
		hasDSC:     false,
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
			this.packages[formatPackageID(packageName.Value, "source")] = struct{}{}
			this.paragraphs = append(this.paragraphs, paragraph)
		} else if _, contains := paragraph.allKeys["Filename"]; contains {
			this.packages[formatPackageID(packageName.Value, architecture.Value)] = struct{}{}
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
	if this.hasDSC && pkg.Architecture() == "source" {
		return false, errors.New("Only a single Debian source package is allowed per manifest.")
	} else if this.contains(pkg) {
		return false, nil
	} else if clone, err := pkg.ToManifest(this.Path()); err != nil {
		return false, err
	} else {
		this.packages[formatPackageID(pkg.Name(), pkg.Architecture())] = struct{}{}
		this.paragraphs = append(this.paragraphs, clone)
		return true, nil
	}
}

func (this *ManifestFile) contains(pkg LocalPackage) bool {
	if _, contains := this.packages[formatPackageID(pkg.Name(), pkg.Architecture())]; contains {
		return true
	} else if pkg.Architecture() == "source" {
		return false
	} else if _, contains := this.packages[formatPackageID(pkg.Name(), "any")]; contains {
		return true
	} else if _, contains := this.packages[formatPackageID(pkg.Name(), "all")]; contains {
		return true
	} else {
		return false
	}
}
func formatPackageID(name, architecture string) string {
	return fmt.Sprintf("%s_%s", name, architecture)
}

func (this *ManifestFile) Bytes() []byte {
	buffer := bytes.NewBuffer([]byte{})
	writer := NewWriter(buffer)
	for _, meta := range this.paragraphs {
		meta.Write(writer)
	}

	return buffer.Bytes()
}
