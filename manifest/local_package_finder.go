package manifest

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/smartystreets/raptr/storage"
)

type LocalPackageFinder struct {
}

func NewLocalPackageFinder(local storage.Storage) *LocalPackageFinder {
	return &LocalPackageFinder{}
}

func (this *LocalPackageFinder) Find(directory string) ([]LocalPackage, error) {
	packages := []LocalPackage{}
	version := ""

	if files, err := ioutil.ReadDir(directory); err != nil {
		return nil, err
	} else {
		for _, file := range files {
			fullPath := path.Join(directory, file.Name())
			if localPackage, err := buildLocalPackage(fullPath); err != nil {
				return nil, err
			} else if localPackage == nil {
				continue
			} else if len(version) > 0 && version != localPackage.Version() {
				return nil, errors.New("All package files must share the same version.")
			} else {
				version = localPackage.Version()
				packages = append(packages, localPackage)
			}
		}
	}

	return packages, nil
}
func buildLocalPackage(fullPath string) (LocalPackage, error) {
	switch strings.ToLower(path.Ext(fullPath)) {
	case "deb":
		return NewPackageFile(fullPath)
	case "dsc":
		return NewSourceFile(fullPath)
	default:
		return nil, nil
	}
}
