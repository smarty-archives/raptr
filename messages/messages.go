package messages

type UploadCommand struct {
	ConfigFile   string // optional
	PackageName  string
	PackagePath  string
	StorageName  string
	Category     string
	Distribution string // optional
}

type LinkCommand struct {
	ConfigFile         string // optional
	PackageName        string
	PackageVersion     string
	StorageName        string
	Category           string
	SourceDistribution string
	TargetDistribution string
}

type UnlinkCommand struct {
	ConfigFile     string // optional
	PackageName    string
	PackageVersion string
	StorageName    string
	Category       string
	Distribution   string
}

type CleanCommand struct {
	ConfigFile  string // optional
	StorageName string
}
