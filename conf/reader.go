package conf

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/smartystreets/raptr/conf/homedir"
	"github.com/smartystreets/raptr/conf/osext"
)

func Read(filenames ...string) (io.Reader, error) {
	return ReadFromLocations(DefaultLocations, filenames...)
}
func ReadFromLocations(locations int, filenames ...string) (io.Reader, error) {
	return ReadSpecific(DefaultLocations, 0, filenames...)
}
func ReadSpecific(locations int, options int, filenames ...string) (io.Reader, error) {
	tried := map[string]struct{}{}
	for _, candidate := range allowedLocations {
		if locations&candidate != candidate {
			continue // they don't want to search that location
		} else if directory, contains := directories[candidate]; !contains {
			continue // location not found in map of allowed directories
		} else if _, contains := tried[directory]; contains {
			continue // we've already looked in this directory
		} else if reader, err := readFiles(options, directory, filenames); os.IsNotExist(err) {
			tried[directory] = struct{}{}
			continue // file not found
		} else if err != nil {
			return nil, err // some unhandled error, e.g. can't read file, permissions issue, etc.
		} else {
			return reader, nil // success
		}
	}

	return nil, os.ErrNotExist
}

func readFiles(options int, directory string, filenames []string) (io.Reader, error) {
	for _, filename := range filenames {
		if reader, err := readFile(path.Join(directory, filename)); os.IsNotExist(err) {
			continue // file doesn't exist, try the next one (if any)
		} else if os.IsPermission(err) && options&SkipPermissionErrors == SkipPermissionErrors {
			continue // permission errors can be skipped
		} else if err != nil {
			return nil, err // general failure, e.g. permissions issues, reading issues, etc.
		} else {
			return reader, nil // success
		}
	}

	return nil, os.ErrNotExist // no files exist
}
func readFile(fullPath string) (io.Reader, error) {
	if handle, err := os.Open(fullPath); err != nil {
		return nil, err // unable to open the file, it may not exist or may have permission issues
	} else if contents, err := ioutil.ReadAll(handle); err != nil {
		handle.Close()
		return nil, err // file can't be read
	} else {
		handle.Close()
		return bytes.NewBuffer(contents), nil // success, convert to io.Reader
	}
}

func init() {
	if workingDirectory, err := os.Getwd(); err == nil {
		addDirectory(Working, workingDirectory)
	}
	if binaryDirectory, err := osext.ExecutableFolder(); err == nil {
		addDirectory(Binary, binaryDirectory)
	}
	if homeDirectory, err := homedir.Dir(); err == nil {
		addDirectory(Home, homeDirectory)
	}

	addDirectory(LocalEtc, "/usr/local/etc")
	addDirectory(Etc, "/etc")
}
func addDirectory(key int, directory string) {
	if len(directory) == 0 {
		return // no directory to add
	} else if absolute, err := filepath.Abs(directory); err != nil {
		return // can't convert to an absolute directory
	} else if _, err := os.Stat(absolute); os.IsNotExist(err) {
		return // it doesn't exist
	} else {
		directories[key] = absolute
	}
}

const (
	Working = 1 << iota
	Binary
	Home
	LocalEtc
	Etc
	DefaultLocations = Working | Binary | Home | LocalEtc | Etc
)
const (
	SkipPermissionErrors = 1 << iota
)

var directories = map[int]string{}
var allowedLocations = []int{Working, Binary, Home, LocalEtc, Etc}
