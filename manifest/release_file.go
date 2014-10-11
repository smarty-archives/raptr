package manifest

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Represents the highest level portion of an APT repository and contains
// checksums of all the various subordinate Packages and Sources files
// for a known set of CPU architectures and software categories
type ReleaseFile struct {
	path          string
	distribution  string
	categories    []string
	architectures []string
	sums          map[string]IndexFile
	items         []IndexFile
}

func NewReleaseFile(distribution string, categories, architectures []string) *ReleaseFile {
	return &ReleaseFile{
		path:          BuildReleaseFilePath(distribution),
		distribution:  distribution,
		categories:    categories,
		architectures: architectures,
		sums:          map[string]IndexFile{},
		items:         []IndexFile{},
	}
}
func BuildReleaseFilePath(distribution string) string {
	return path.Join("/dists/", distribution, "Release")
}

func (this *ReleaseFile) Add(index IndexFile) bool {
	added := false
	if _, contains := this.sums[index.Path()]; !contains {
		added = true
		this.sums[index.Path()] = index
		this.items = append(this.items, index)
	}
	return added
}

func (this *ReleaseFile) Parse(reader io.Reader) error {
	return nil // TODO
}

func (this *ReleaseFile) Bytes() []byte {
	paragraph := NewParagraph()

	addLine(paragraph, "Architectures", strings.Join(this.architectures, " "))
	addLine(paragraph, "Components", strings.Join(this.categories, " "))
	addLine(paragraph, "Date", time.Now().UTC().Format(time.RFC1123))
	addLine(paragraph, "Description", "none")
	addLine(paragraph, "Origin", "raptr")
	addLine(paragraph, "Suite", this.distribution)

	checksums := []Checksum{}
	for _, item := range this.items {
		checksum, _ := ComputeChecksums(bytes.NewBuffer(item.Bytes()))
		checksums = append(checksums, checksum)
	}
	addLine(paragraph, "MD5Sum", "")
	for i, item := range this.items {
		this.addHashLine(paragraph, item, checksums[i].MD5)
	}
	addLine(paragraph, "SHA1Sum", "")
	for i, item := range this.items {
		this.addHashLine(paragraph, item, checksums[i].SHA1)
	}
	addLine(paragraph, "SHA256Sum", "")
	for i, item := range this.items {
		this.addHashLine(paragraph, item, checksums[i].SHA256)
	}
	addLine(paragraph, "SHA512Sum", "")
	for i, item := range this.items {
		this.addHashLine(paragraph, item, checksums[i].SHA512)
	}

	return serializeParagraphs([]*Paragraph{paragraph})
}
func (this *ReleaseFile) addHashLine(paragraph *Paragraph, item IndexFile, checksum []byte) {
	basepath := filepath.Dir(this.Path())
	relativePath, _ := filepath.Rel(basepath, item.Path())
	line := fmt.Sprintf("%x %16d %s", checksum, len(item.Bytes()), relativePath)
	addLine(paragraph, "", line)
}

func (this *ReleaseFile) Path() string {
	return this.path
}
