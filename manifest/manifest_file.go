package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
)

type ManifestFile struct {
	category      string
	bundle        string
	version       string
	hasDSC        bool
	paragraphs    []*Paragraph
	packages      map[string]struct{}
	architectures map[string][]*Paragraph
}

func NewManifestFile(category, bundle, version string) *ManifestFile {
	return &ManifestFile{
		category:      category,
		bundle:        bundle,
		version:       version,
		hasDSC:        false,
		paragraphs:    []*Paragraph{},
		packages:      map[string]struct{}{},
		architectures: map[string][]*Paragraph{},
	}
}
func ParseManifest(reader io.Reader, category, bundle, version string) (*ManifestFile, error) {
	this := NewManifestFile(category, bundle, version)
	// gzipReader, err := gzip.NewReader(reader)
	// paragraphReader := NewReader(gzipReader)
	// if err != nil {
	// 	return nil, err
	// }
	paragraphReader := NewReader(reader)

	for {
		if paragraph, err := ReadParagraph(paragraphReader); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if packageName, contains := paragraph.allKeys["Package"]; !contains {
			return nil, errors.New("Malformed manifest file, missing Package element.")
		} else if architecture, contains := paragraph.allKeys["Architecture"]; !contains {
			return nil, errors.New("Malformed manifest file, missing Architecture element.")
		} else if pkgVersion, contains := paragraph.allKeys["Version"]; !contains {
			return nil, errors.New("Malformed manifest file, missing Version element.")
		} else if version != pkgVersion.Value {
			return nil, errors.New("The package version differs from the provided manifest version.")
		} else if _, contains := paragraph.allKeys["Files"]; contains {
			this.packages[formatPackageID(packageName.Value, "source")] = struct{}{}
			this.architectures["source"] = append(this.architectures["source"], paragraph)
			this.paragraphs = append(this.paragraphs, paragraph)
		} else if _, contains := paragraph.allKeys["Filename"]; contains {
			this.packages[formatPackageID(packageName.Value, architecture.Value)] = struct{}{}
			this.architectures[architecture.Value] = append(this.architectures[architecture.Value], paragraph)
			this.paragraphs = append(this.paragraphs, paragraph)
		}
	}

	return this, nil
}
func BuildPath(category, bundle, version string) string {
	return path.Join("/pool/", category, bundle[0:1], bundle, version, "manifest")
}

func (this *ManifestFile) Architectures() []string {
	items := []string{}
	for key, _ := range this.architectures {
		items = append(items, key)
	}
	return items
}

func (this *ManifestFile) Path() string {
	return path.Dir(BuildPath(this.category, this.bundle, this.version))
}
func (this *ManifestFile) Add(pkg LocalPackage) (bool, error) {
	if this.hasDSC && pkg.Architecture() == "source" {
		return false, errors.New("Only a single Debian source package is allowed per manifest.")
	} else if pkg.Version() != this.version {
		return false, errors.New("The package to be added does not match the manifest version.")
	} else if this.contains(pkg) {
		return false, nil
	} else if clone, err := pkg.ToManifest(this.Path()); err != nil {
		return false, err
	} else {
		this.packages[formatPackageID(pkg.Name(), pkg.Architecture())] = struct{}{}
		this.architectures[pkg.Architecture()] = append(this.architectures[pkg.Architecture()], clone)
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
	return serializeParagraphs(this.paragraphs)
}
func serializeParagraphs(paragraphs []*Paragraph) []byte {
	buffer := bytes.NewBuffer([]byte{})
	// gzipWriter, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)

	// writer := NewWriter(gzipWriter)
	writer := NewWriter(buffer)
	for _, paragraph := range paragraphs {
		paragraph.Write(writer)
	}

	// gzipWriter.Close()
	return buffer.Bytes()
}
