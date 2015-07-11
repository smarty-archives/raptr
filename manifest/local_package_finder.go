package manifest

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type LocalPackageFinder struct{}

func NewLocalPackageFinder() *LocalPackageFinder {
	return &LocalPackageFinder{}
}

func (this *LocalPackageFinder) Find(directory string) ([]LocalPackage, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	version := ""
	packages := []LocalPackage{}

	for item := range discoverAllPackages(directory, files) {
		if item == nil {
			continue
		} else if err, ok := item.(error); ok {
			return nil, err
		} else if pkg, ok := item.(LocalPackage); !ok {
			continue
		} else if len(version) > 0 && version != pkg.Version() {
			return nil, errors.New("All package files must share the same version.")
		} else {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

func discoverAllPackages(directory string, files []os.FileInfo) chan interface{} {
	result := make(chan interface{}, 256)
	waiter := sync.WaitGroup{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		waiter.Add(1)
		go func(fullPath string) {
			result <- discoverPackage(fullPath)
			waiter.Done()
		}(path.Join(directory, file.Name()))
	}
	waiter.Wait()
	close(result)
	return result
}

func discoverPackage(fullPath string) interface{} {
	if pkg, err := buildLocalPackage(fullPath); err != nil {
		return err
	} else {
		return pkg
	}
}

func buildLocalPackage(fullPath string) (LocalPackage, error) {
	switch strings.ToLower(path.Ext(fullPath)) {
	case ".deb":
		return NewPackageFile(fullPath)
	case ".dsc":
		return NewSourceFile(fullPath)
	default:
		return nil, nil
	}
}
