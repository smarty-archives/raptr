package manifest

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
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
	} else if info, err := os.Stat(fullPath); err != nil {
		return nil, err
	} else if handle, err := os.Open(fullPath); err != nil {
		return nil, err
	} else if computed, err := ComputeChecksums(handle); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else if paragraph, err := ReadParagraph(NewReader(handle)); err != nil {
		handle.Close()
		return nil, err
	} else if _, err := handle.Seek(0, 0); err != nil {
		handle.Close()
		return nil, err
	} else if files, err := readSourcePackageFiles(fullPath, paragraph); err != nil {
		handle.Close()
		return nil, err
	} else if len(files) == 0 {
		handle.Close()
		return nil, errors.New("Debian source code packages does not contain any files.")
	} else {
		// TODO: ensure that contents of file agree with filename scheme
		file := LocalPackageFile{
			Name:      strings.ToLower(path.Base(fullPath)),
			Length:    uint64(info.Size()),
			Checksums: computed,
			Contents:  handle,
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
		} else if packageArchive := path.Join(path.Dir(fullPath), filename); false {
			return nil, nil
		} else if info, err := os.Stat(packageArchive); err != nil {
			return nil, err
		} else if handle, err := os.Open(packageArchive); err != nil {
			return nil, err
		} else if computed, err := ComputeChecksums(handle); err != nil {
			handle.Close()
			return nil, err
		} else if bytes.Compare(computed.MD5, parsedMD5) != 0 {
			handle.Close()
			return nil, errors.New("File contents do not match line item in dsc file.")
		} else if _, err := handle.Seek(0, 0); err != nil {
			handle.Close()
			return nil, err
		} else {
			files = append(files, LocalPackageFile{
				Name:      filename,
				Length:    uint64(info.Size()),
				Checksums: computed,
				Contents:  handle,
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
	if len(split) < 2 {
		return "", ""
	} else {
		return split[0], split[len(split)-1]
	}
}
func (this *SourceFile) ToManifest(poolDirectory string) (*Paragraph, error) {

	if clone := this.cloneWithoutFiles(); !clone.RenameKey("Source", "Package") {
		return nil, errors.New("Unable to rename desired key")
	} else if !addLine(clone, "Directory", poolDirectory) {
		return nil, errors.New("Unable to add a line to the debian control file")
	} else {
		addLine(clone, "Checksums-Sha1", "")
		for _, file := range this.Files() {
			addLine(clone, "", fmt.Sprintf("%x %d %s", file.Checksums.SHA1, file.Length, file.Name))
		}
		addLine(clone, "Checksums-Sha256", "")
		for _, file := range this.Files() {
			addLine(clone, "", fmt.Sprintf("%x %d %s", file.Checksums.SHA256, file.Length, file.Name))
		}
		addLine(clone, "Files", "")
		for _, file := range this.Files() {
			addLine(clone, "", fmt.Sprintf("%x %d %s", file.Checksums.MD5, file.Length, file.Name))
		}

		return clone, nil
	}
}
func addLine(meta *Paragraph, key, value string) bool {
	if line, err := NewLine(key, value); err != nil {
		return false
	} else if err := meta.Add(line, false); err != nil {
		return false
	} else {
		return true
	}
}
func (this *SourceFile) cloneWithoutFiles() *Paragraph {
	clone := NewParagraph()
	skip := false
	checksumPrefix := normalizeKey("Checksums-")
	filesKey := normalizeKey("Files")

	for _, item := range this.paragraph.items {
		if skip && len(item.Key) == 0 {
			continue // skip any value-only lines when we're in skip mode
		} else if skip = item.Key == filesKey; skip {
			continue // skip the "Files:" section
		} else if skip = strings.HasPrefix(item.Key, checksumPrefix); skip {
			continue // skip the "Checksum-*" section
		}

		clone.allKeys[item.Key] = item
		clone.items = append(clone.items, item)
		clone.orderedKeys = append(clone.orderedKeys, item.Key)
	}

	return clone
}

func (this *SourceFile) Name() string              { return this.name }
func (this *SourceFile) Version() string           { return this.version }
func (this *SourceFile) Architecture() string      { return "source" }
func (this *SourceFile) Files() []LocalPackageFile { return this.files }
