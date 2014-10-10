package manifest

// Represents the highest level portion of an APT repository and contains
// checksums of all the various subordinate Packages and Sources files
// for a known set of CPU architectures and software categories
// NOTE: it may be that this is a write-only file (depending upon the application logic)
// and concurrency-related issues
type ReleaseFile struct {
	path string
}

func NewReleaseFile() *ReleaseFile {
	return &ReleaseFile{path: BuildReleaseFilePath()}
}
func BuildReleaseFilePath() string {
	return "/Release.gz"
}

func (this *ReleaseFile) Bytes() []byte {
	// TODO: implements a gzip writer to compress stuff
	return []byte{}
}

func (this *ReleaseFile) Path() string {
	return this.path
}
