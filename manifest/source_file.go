package manifest

import (
	"bytes"
	"encoding/hex"
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
	} else if files, err := readSourcePackageFiles(fullPath, paragraph); err != nil {
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
func readSourcePackageFiles(fullPath string, paragraph *Paragraph) ([]LocalPackageFile, error) {
	files := []LocalPackageFile{}

	for _, line := range getSourcePackageFilenameLineValues(paragraph) {
		md5hash, filename := parseFileLine(line)
		if len(md5hash) == 0 || len(filename) == 0 {
			return nil, errors.New("Unable to parse line")
		} else if parsedMD5, err := hex.DecodeString(md5hash); err != nil {
			return nil, err
		} else if handle, err := os.Open(path.Join(path.Dir(fullPath), filename)); err != nil {
			return nil, err
		} else if computed, err := computeMD5(handle); err != nil {
			handle.Close()
			return nil, err
		} else if bytes.Compare(computed, parsedMD5) != 0 {
			handle.Close()
			return nil, errors.New("File contents do not match line item in dsc file.")
		} else if _, err := handle.Seek(0, 0); err != nil {
			handle.Close()
			return nil, err
		} else {
			files = append(files, LocalPackageFile{
				Name:     filename,
				MD5:      computed,
				Contents: handle,
			})
		}
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
func parseFileLine(value string) (string, string) {
	split := strings.Split(value, " ")
	if len(split) > 2 {
		return "", ""
	} else {
		return split[0], split[len(split)-1]
	}
}

func (this *SourceFile) Name() string              { return this.name }
func (this *SourceFile) Version() string           { return this.version }
func (this *SourceFile) Architecture() string      { return "source" }
func (this *SourceFile) Files() []LocalPackageFile { return this.files }
