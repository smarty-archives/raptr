package manifest

// Represents the inner contents of a compiled, debian archive
// NOTE: we will only ever read these files
type PackageFile struct{}

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
}
