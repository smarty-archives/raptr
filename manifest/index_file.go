package manifest

type IndexFile interface {
	Path() string
	Add(*ManifestFile)
}
