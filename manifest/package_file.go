package manifest

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/blakesmith/ar"
)

// Represents the inner contents of a compiled, debian archive
// NOTE: we will only ever read these files
type PackageFile struct {
	name         string
	version      string
	architecture string
	paragraph    *Paragraph
	file         LocalPackageFile
}

func NewPackageFile(fullPath string) (*PackageFile, error) {
	if meta := ParseFilename(fullPath); meta == nil {
		return nil, errors.New("The file provided is not a debian binary package.")
	} else if info, err := os.Stat(fullPath); err != nil {
		return nil, err
	} else if handle, err := os.Open(fullPath); err != nil {
		return nil, err
	} else if computed, err := computeChecksums(fullPath, handle); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else if paragraph, err := extractManifest(fullPath, handle); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else {
		// TODO: ensure that contents of internal control file agree with filename scheme
		return &PackageFile{
			name:         meta.Name,
			version:      paragraph.Version(),
			architecture: meta.Architecture,
			paragraph:    paragraph,
			file: LocalPackageFile{
				Name:      path.Base(fullPath),
				Length:    uint64(info.Size()),
				Checksums: computed,
				Contents:  handle,
			},
		}, nil
	}
}

func computeChecksums(fullPath string, reader io.Reader) (Checksum, error) {
	log.Println("[INFO] Computing checksums for", path.Base(fullPath))
	defer log.Println("[INFO] Finished computing checksums for", path.Base(fullPath))
	return ComputeChecksums(reader)
}

func extractManifest(fullPath string, reader io.Reader) (*Paragraph, error) {
	archiveReader := ar.NewReader(reader)

	log.Println("[INFO] Extracting debian/control file from", path.Base(fullPath))

	for {
		if archiveHeader, err := archiveReader.Next(); err != nil {
			return nil, err
		} else if path.Base(archiveHeader.Name) != "control.tar.gz" {
			continue
		} else if gzipReader, err := gzip.NewReader(archiveReader); err != nil {
			return nil, err
		} else if tarReader := tar.NewReader(gzipReader); false {
			continue
		} else {
			for {
				if fileHeader, err := tarReader.Next(); err != nil {
					return nil, err
				} else if path.Base(fileHeader.Name) != "control" {
					continue
				} else if paragraph, err := ReadParagraph(NewReader(tarReader)); err != nil {
					return nil, err
				} else {
					return paragraph, nil
				}
			}
		}
	}
}

func (this *PackageFile) ToManifest(poolDirectory string) (*Paragraph, error) {
	clone := NewParagraph()
	added := false

	for _, item := range this.paragraph.items {
		if item.Key == normalizeKey("Depends") {
			addLine(clone, item.Key, item.Value)
			this.addChecksumLines(clone, poolDirectory)
			added = true
		} else if item.Key == normalizeKey("Section") && !added {
			this.addChecksumLines(clone, poolDirectory)
			addLine(clone, item.Key, item.Value)
			added = true
		} else if item.Key == normalizeKey("Description") && !added {
			this.addChecksumLines(clone, poolDirectory)
			addLine(clone, item.Key, item.Value)
			added = true
		} else if item.Key == normalizeKey("Version") {
			addLine(clone, item.Key, this.version)
		} else {
			addLine(clone, item.Key, item.Value)
		}
	}

	return clone, nil
}
func (this *PackageFile) addChecksumLines(clone *Paragraph, directory string) {
	addLine(clone, "Filename", path.Join(directory[1:], this.file.Name))
	addLine(clone, "Size", fmt.Sprintf("%d", this.file.Length))
	addLine(clone, "MD5sum", fmt.Sprintf("%x", this.file.Checksums.MD5))
	addLine(clone, "SHA1", fmt.Sprintf("%x", this.file.Checksums.SHA1))
	addLine(clone, "SHA256", fmt.Sprintf("%x", this.file.Checksums.SHA256))
}

func (this *PackageFile) Name() string              { return this.name }
func (this *PackageFile) Version() string           { return this.version }
func (this *PackageFile) Architecture() string      { return this.architecture }
func (this *PackageFile) Files() []LocalPackageFile { return []LocalPackageFile{this.file} }
