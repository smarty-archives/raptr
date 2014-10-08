package manifest

import "os"

// Represents the inner contents of a compiled, debian archive
// NOTE: we will only ever read these files
type PackageFile struct {
	name         string
	version      string
	architecture string
	handle       *os.File
}

func NewPackageFile(fullPath string) (*PackageFile, error) {
	return nil, nil
	// name := path.Base(fullpath)
	// split := strings.Split(name, "_")
	// if len(split) != 3 {

	// }

	// if _, err := os.Stat(); err != nil {
	// 	return nil, err
	// }

	// /path/to/file/raptr2_1.0.7_amd64.deb
	// TODO: parse name--it must be well formed to be a deb package file
}

func (this *PackageFile) Name() string {
	return ""
}
func (this *PackageFile) Version() string {
	return ""
}
func (this *PackageFile) Architecture() string {
	return ""
}
func (this *PackageFile) Files() []LocalPackageFile {
	return nil
	// return []LocalPackageFile{
	// 	LocalPackage{},
	// }
}
