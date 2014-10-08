package manifest

// Represents a "dsc" file which results from building a debian source code package
// NOTE: we will only ever read these files
type SourceFile struct{}

func NewSourceFile(fullPath string) (*SourceFile, error) {
	return nil, nil
}

func (this *SourceFile) Name() string {
	return ""
}
func (this *SourceFile) Version() string {
	return ""
}
func (this *SourceFile) Architecture() string {
	return ""
}
func (this *SourceFile) Files() []LocalPackageFile {
	return nil
}
