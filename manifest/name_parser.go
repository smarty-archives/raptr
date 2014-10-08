package manifest

type ParsedName struct {
	Name         string
	Version      string
	Architecture string
}

func ParseName(name string) ParsedName {
	return ParsedName{}
}
