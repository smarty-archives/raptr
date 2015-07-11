package manifest

import (
	"errors"
	"io/ioutil"
	"log"
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
			version = pkg.Version()
		}
	}

	return packages, nil
}

func discoverAllPackages(directory string, files []os.FileInfo) chan interface{} {
	log.Println("Discovering all local packages.")
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
	log.Println("Waiting for local file parsing...")
	waiter.Wait()
	log.Println("Local file parsing complete.")
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
		log.Println("Parsing binary package:", fullPath)
		return NewPackageFile(fullPath)
	case ".dsc":
		log.Println("Parsing source package:", fullPath)
		return NewSourceFile(fullPath)
	default:
		return nil, nil
	}
}
