package manifest

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Represents the highest level portion of an APT repository and contains
// checksums of all the various subordinate Packages and Sources files
// for a known set of CPU architectures and software categories
type ReleaseFile struct {
	filePath      string
	baseDirectory string
	distribution  string
	categories    []string
	architectures []string
	items         map[string]func() ReleaseItem
	cached        []byte
}
type ReleaseItem struct {
	RelativePath string
	Length       uint64
	Checksums    map[string][]byte
}

func NewReleaseFile(distribution string, categories, architectures []string) *ReleaseFile {
	filePath := BuildReleaseFilePath(distribution)
	return &ReleaseFile{
		filePath:      filePath,
		baseDirectory: filepath.Dir(filePath),
		distribution:  distribution,
		categories:    categories,
		architectures: architectures,
		items:         map[string]func() ReleaseItem{},
	}
}
func BuildReleaseFilePath(distribution string) string {
	return path.Join("/dists/", distribution, "Release")
}

func (this *ReleaseFile) Add(index IndexFile) bool {
	this.cached = nil
	basepath := filepath.Dir(this.Path())
	relativePath, _ := filepath.Rel(basepath, index.Path())
	this.items[relativePath] = func() ReleaseItem {
		return this.translateIndexFile(relativePath, index)
	}
	return true
}
func (this *ReleaseFile) translateIndexFile(relativePath string, file IndexFile) ReleaseItem {
	checksum, _ := ComputeChecksums(bytes.NewBuffer(file.Bytes()))
	checksums := map[string][]byte{}
	checksums["MD5"] = checksum.MD5
	checksums["SHA1"] = checksum.SHA1
	checksums["SHA256"] = checksum.SHA256
	return ReleaseItem{
		RelativePath: relativePath,
		Length:       uint64(len(file.Bytes())),
		Checksums:    checksums,
	}
}

func (this *ReleaseFile) Parse(reader io.Reader) error {
	this.cached = nil
	this.items = map[string]func() ReleaseItem{}

	paragraph, err := ReadParagraph(NewReader(reader))
	if err != nil {
		return err
	}

	parsed := map[string]ReleaseItem{}
	hashType := ""
	for _, item := range paragraph.items {
		if strings.HasSuffix(item.Key, "Sum") {
			hashType = item.Key[0 : len(item.Key)-3]
		} else if item.Key == "SHA1" || item.Key == "SHA256" {
			hashType = item.Key
		} else if len(hashType) == 0 {
			continue
		} else if err := parseReleaseItem(hashType, item.Value, parsed); err != nil {
			return err
		}
	}

	for _, item := range parsed {
		unique := item
		this.items[unique.RelativePath] = func() ReleaseItem {
			return unique
		}
	}

	return nil
}
func parseReleaseItem(hashType, unparsed string, parsed map[string]ReleaseItem) error {
	unparsed = strings.TrimSpace(unparsed)
	indexOfWhitespace := strings.Index(unparsed, " ")
	if indexOfWhitespace == -1 {
		return errors.New("Malformed line--missing hash")
	}

	computed, err := hex.DecodeString(unparsed[0:indexOfWhitespace])
	if err != nil {
		return errors.New("Malformed line--bad hash")
	}
	unparsed = strings.TrimSpace(unparsed[indexOfWhitespace+1:])

	indexOfWhitspace := strings.LastIndex(unparsed, " ")
	if indexOfWhitspace == -1 {
		return errors.New("Malformed line--missing filename")
	}

	relativePath := unparsed[indexOfWhitspace+1:]
	unparsed = strings.TrimSpace(unparsed[0:indexOfWhitspace])

	length, err := strconv.ParseUint(unparsed, 10, 64)
	if err != nil {
		return errors.New("Malformed line--missing length")
	}

	item, contains := parsed[relativePath]
	if !contains {
		item.RelativePath = relativePath
		item.Length = length
		item.Checksums = map[string][]byte{}
		parsed[relativePath] = item
	}

	item.Checksums[hashType] = computed
	return nil
}

func (this *ReleaseFile) Bytes() []byte {
	if this.cached != nil {
		return this.cached
	}

	paragraph := NewParagraph()

	addLine(paragraph, "Architectures", strings.Join(this.architectures, " "))
	addLine(paragraph, "Components", strings.Join(this.categories, " "))
	addLine(paragraph, "Date", time.Now().UTC().Format(time.RFC1123))
	addLine(paragraph, "Description", "none")
	addLine(paragraph, "Origin", "raptr")
	addLine(paragraph, "Suite", this.distribution)

	keys := []string{}
	for key, _ := range this.items {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	releaseItems := []ReleaseItem{}
	for _, key := range keys {
		releaseItems = append(releaseItems, this.items[key]())
	}

	for _, hashType := range []string{"MD5", "SHA1", "SHA256"} {
		hashHeader := hashType
		if hashType == "MD5" {
			hashHeader = hashType + "Sum" // MD5Sum vs SHA1, SHA256, etc.
		}

		addLine(paragraph, hashHeader, "")
		for _, item := range releaseItems {
			value := fmt.Sprintf("%x %16d %s", item.Checksums[hashType], item.Length, item.RelativePath)
			addLine(paragraph, "", value)
		}
	}

	this.cached = serializeParagraphs([]*Paragraph{paragraph})
	return this.cached
}

func (this *ReleaseFile) Path() string {
	return this.filePath
}
