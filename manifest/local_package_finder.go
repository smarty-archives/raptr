package manifest

import (
	"errors"
	"io/ioutil"
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

	packages := []LocalPackage{}
	waiter, mutex := sync.WaitGroup{}, sync.Mutex{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		waiter.Add(1)
		go func(fileName string) {
			pkg, err := buildLocalPackage(directory, fileName)
			if pkg != nil {
				mutex.Lock()

				expected := pkg.Version()
				for _, item := range packages {
					if expected != item.Version() {
						err = errors.New("All package files must share the same version.")
					}
				}

				packages = append(packages, pkg)
				mutex.Unlock()
			}
			if err != nil {
				mutex.Lock()
				err = err
				mutex.Unlock()
			}
			waiter.Done()
		}(file.Name())
	}

	waiter.Wait()

	if err != nil {
		return nil, err
	} else {
		return packages, nil
	}
}

func buildLocalPackage(directory, filename string) (LocalPackage, error) {
	fullPath := path.Join(directory, filename)
	switch strings.ToLower(path.Ext(fullPath)) {
	case ".deb":
		return NewPackageFile(fullPath)
	case ".dsc":
		return NewSourceFile(fullPath)
	default:
		return nil, nil
	}
}
