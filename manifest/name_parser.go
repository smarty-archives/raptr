package manifest

import (
	"path"
	"strings"
)

type ParsedName struct {
	Name         string
	Version      string
	Architecture string
	Container    string // if it exists
}

func ParseFilename(fullPath string) *ParsedName {
	filename := strings.ToLower(path.Base(fullPath))
	extension := path.Ext(filename)
	filename = strings.TrimSuffix(filename, extension)
	if len(extension) > 0 {
		extension = extension[1:]
	}
	parts := strings.Split(filename, "_")

	switch extension {
	case "deb":
		if len(parts) != 3 {
			return nil
		}
		return &ParsedName{
			Name:         parts[0],
			Version:      parts[1],
			Architecture: parts[2],
			Container:    extension,
		}
	case "dsc":
		if len(parts) != 2 {
			return nil
		}
		return &ParsedName{
			Name:         parts[0],
			Version:      parts[1],
			Architecture: "source",
			Container:    extension,
		}
	default:
		return nil
	}
}
