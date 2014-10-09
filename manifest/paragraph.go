package manifest

import (
	"errors"
	"io"
	"strings"
)

type Paragraph struct {
	allKeys     map[string]*LineItem
	orderedKeys []string
	items       []*LineItem
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		allKeys:     map[string]*LineItem{},
		orderedKeys: []string{},
		items:       []*LineItem{},
	}
}

func ReadParagraph(reader *Reader) (*Paragraph, error) {
	this := NewParagraph()

	for {
		if item, err := reader.Read(); len(this.items) > 0 && err == io.EOF {
			return this, nil
		} else if err != nil {
			return nil, err
		} else if item.Type == separator && len(this.items) == 0 {
			continue // keep reading until we find something
		} else if item.Type == separator {
			break
		} else if err := this.Add(item, false); err != nil {
			return nil, err
		}
	}

	return this, nil
}

func (this *Paragraph) Add(item *LineItem, overwrite bool) error {
	if item == nil {
		return nil
	} else if item.Type == separator {
		return nil
	} else if item.Type == comment {
		this.items = append(this.items, item)
	} else if item.Type == valueOnly {
		this.items = append(this.items, item)
	} else if normalized := normalizeKey(item.Key); len(normalized) == 0 {
		return nil
	} else if _, contains := this.allKeys[normalized]; contains && !overwrite {
		return errors.New("The paragraph already contains the specified key")
	} else if !contains {
		this.allKeys[normalized] = item
		this.orderedKeys = append(this.orderedKeys, item.Key)
		this.items = append(this.items, item)
		return nil
	} else {
		// overwrite
		this.allKeys[normalized] = item
		for i, existing := range this.items {
			if normalized == strings.ToTitle(existing.Key) {
				this.items[i] = item
				break
			}
		}
	}

	return nil
}
func (this *Paragraph) RenameKey(oldKey, newKey string) bool {
	oldKey = normalizeKey(oldKey)
	newKey = normalizeKey(newKey)

	if len(oldKey) == 0 || len(newKey) == 0 {
		return false
	} else if _, contains := this.allKeys[newKey]; contains {
		return false // can't replace if the new one already exists
	} else if item, contains := this.allKeys[oldKey]; !contains {
		return false
	} else {
		delete(this.allKeys, oldKey)
		this.allKeys[newKey] = item
		item.Key = newKey

		for i, item := range this.orderedKeys {
			if item == oldKey {
				this.orderedKeys[i] = newKey
				break
			}
		}

		return true
	}
}

func normalizeKey(key string) string {
	if len(key) == 0 {
		return key
	} else if key = strings.TrimSpace(strings.ToLower(key)); len(key) == 0 {
		return key
	} else if key == "md5sum" {
		return "MD5Sum"
	} else if key == "sha1sum" {
		return "SHA1Sum"
	} else if key == "sha256sum" {
		return "SHA256Sum"
	} else if key == "sha512sum" {
		return "SHA512Sum"
	} else {
		return strings.Title(key)
	}
}

func (this *Paragraph) Write(writer *Writer) error {
	if len(this.items) == 0 {
		return nil
	}

	for _, item := range this.items {
		if err := writer.Write(item); err != nil {
			return err
		}
	}

	return writer.Write(&LineItem{Type: separator})
}
