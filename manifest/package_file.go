package manifest

import (
	"errors"
	"os"
	"path"
	"strings"
)

// Represents the inner contents of a compiled, debian archive
// NOTE: we will only ever read these files
type PackageFile struct {
	name         string
	version      string
	architecture string
	file         LocalPackageFile
}

func NewPackageFile(fullPath string) (*PackageFile, error) {
	if meta := ParseFilename(fullPath); meta == nil {
		return nil, errors.New("The file provided is not a debian binary package.")
	} else if handle, err := os.Open(fullPath); err != nil {
		return nil, err
	} else if computed, err := ComputeChecksums(handle); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else {
		// TODO: *open* the debian file and read the manifest/control file
		return &PackageFile{
			name:         meta.Name,
			version:      meta.Version,
			architecture: meta.Architecture,
			file: LocalPackageFile{
				Name:      strings.ToLower(path.Base(fullPath)),
				Contents:  handle,
				Checksums: computed,
			},
		}, nil
	}
}

func (this *PackageFile) Name() string              { return this.name }
func (this *PackageFile) Version() string           { return this.version }
func (this *PackageFile) Architecture() string      { return this.architecture }
func (this *PackageFile) Files() []LocalPackageFile { return []LocalPackageFile{this.file} }
