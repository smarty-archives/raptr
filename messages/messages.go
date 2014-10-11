package messages

type UploadCommand struct {
	PackageName  string
	PackagePath  string
	StorageName  string
	Category     string
	Distribution string // optional
}

type LinkCommand struct {
	PackageName    string
	PackageVersion string
	StorageName    string
	Category       string
	Distribution   string
}

type UnlinkCommand struct {
	PackageName    string
	PackageVersion string
	StorageName    string
	Category       string
	Distribution   string
}

type CleanCommand struct {
	StorageName string
}
