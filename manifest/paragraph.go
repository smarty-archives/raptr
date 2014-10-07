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
		if item, err := reader.Read(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
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
	} else if title := strings.ToTitle(item.Key); len(title) == 0 {
		return nil
	} else if _, contains := this.allKeys[title]; contains && !overwrite {
		return errors.New("The paragraph already contains the specified key")
	} else if !contains {
		this.allKeys[title] = item
		this.orderedKeys = append(this.orderedKeys, item.Key)
		this.items = append(this.items, item)
		return nil
	} else {
		// overwrite
		this.allKeys[title] = item
		for i, existing := range this.items {
			if title == strings.ToTitle(existing.Key) {
				this.items[i] = item
				break
			}
		}
	}

	return nil
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
