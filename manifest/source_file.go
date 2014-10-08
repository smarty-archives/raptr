package manifest

import (
	"errors"
	"os"
	"path"
	"strings"
)

// Represents a "dsc" file which results from building a debian source code package
// NOTE: we will only ever read these files
type SourceFile struct {
	name      string
	version   string
	files     []LocalPackageFile
	paragraph *Paragraph
}

func NewSourceFile(fullPath string) (*SourceFile, error) {
	if meta := ParseFilename(fullPath); meta == nil {
		return nil, errors.New("The file provided is not a debian binary package.")
	} else if handle, err := os.Open(fullPath); err != nil {
		return nil, err
	} else if computed, err := computeMD5(handle); err != nil {
		return nil, err
	} else if paragraph, err := ReadParagraph(NewReader(handle)); err != nil {
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		return nil, err
	} else if files, err := readSourcePackageFiles(paragraph); err != nil {
		return nil, err
	} else if len(files) == 0 {
		return nil, errors.New("Debian source code packages does not contain any files.")
	} else {
		file := LocalPackageFile{
			Name:     strings.ToLower(path.Base(fullPath)),
			Contents: handle,
			MD5:      computed,
		}
		files = append([]LocalPackageFile{file}, files...)

		return &SourceFile{
			name:      meta.Name,
			version:   meta.Version,
			files:     files,
			paragraph: paragraph,
		}, nil
	}
}
func readSourcePackageFiles(paragraph *Paragraph) ([]LocalPackageFile, error) {
	files := []LocalPackageFile{}

	for _, line := range getSourcePackageFilenameLineValues(paragraph) {
		line = line
		// split the line to figure out the name of the file, the size, and the MD5 hash
		// the filename should not contain any paths and should be alongside the dsc
		// open the file, it should exist
		// finally compute the MD5 of the file and compare to the line item
		// append to files package
		// if there are any problems above, fail and return error
	}

	return files, nil
}
func getSourcePackageFilenameLineValues(paragraph *Paragraph) []string {
	fileLines := []string{}
	foundFilesKey := false
	for _, item := range paragraph.items {
		if foundFilesKey && len(item.Key) > 0 {
			break
		} else if foundFilesKey {
			fileLines = append(fileLines, item.Value)
		} else {
			foundFilesKey = item.Key == "Files"
		}
	}

	return fileLines
}

func (this *SourceFile) Name() string              { return this.name }
func (this *SourceFile) Version() string           { return this.version }
func (this *SourceFile) Architecture() string      { return "source" }
func (this *SourceFile) Files() []LocalPackageFile { return this.files }
