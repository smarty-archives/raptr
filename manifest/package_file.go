package manifest

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

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
	} else if computed, err := ComputeChecksums(handle); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else if paragraph, err := extractManifest(handle); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else {
		// TODO: ensure that contents of internal control file agree with filename scheme
		return &PackageFile{
			name:         meta.Name,
			version:      meta.Version,
			architecture: meta.Architecture,
			paragraph:    paragraph,
			file: LocalPackageFile{
				Name:      strings.ToLower(path.Base(fullPath)),
				Length:    uint64(info.Size()),
				Checksums: computed,
				Contents:  handle,
			},
		}, nil
	}
}
func extractManifest(reader io.Reader) (*Paragraph, error) {
	archiveReader := ar.NewReader(reader)

	for {
		if archiveHeader, err := archiveReader.Next(); err != nil {
			return nil, err
		} else if strings.ToLower(path.Base(archiveHeader.Name)) != "control.tar.gz" {
			continue
		} else if gzipReader, err := gzip.NewReader(archiveReader); err != nil {
			return nil, err
		} else if tarReader := tar.NewReader(gzipReader); false {
			continue
		} else {
			for {
				if fileHeader, err := tarReader.Next(); err != nil {
					return nil, err
				} else if strings.ToLower(path.Base(fileHeader.Name)) != "control" {
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
	clone := this.paragraph.CloneWithoutFiles()

	// TODO: order of adding lines? perhaps we clone again but watch
	// for a certain attribute and insert these before it?
	addLine(clone, "Filename", path.Join(poolDirectory, this.file.Name))
	addLine(clone, "Size", fmt.Sprintf("%d", this.file.Length))
	addLine(clone, "MD5sum", fmt.Sprintf("%x", this.file.Checksums.MD5))
	addLine(clone, "SHA1", fmt.Sprintf("%x", this.file.Checksums.SHA1))
	addLine(clone, "SHA256", fmt.Sprintf("%x", this.file.Checksums.SHA256))

	return clone, nil
}

func (this *PackageFile) Name() string              { return this.name }
func (this *PackageFile) Version() string           { return this.version }
func (this *PackageFile) Architecture() string      { return this.architecture }
func (this *PackageFile) Files() []LocalPackageFile { return []LocalPackageFile{this.file} }
