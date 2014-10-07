package manifest

import (
	"errors"
	"io"
	"strings"
)

type Paragraph struct {
	keys  map[string]*LineItem
	items []*LineItem
}

func ReadParagraph(reader *Reader) (*Paragraph, error) {
	keys := map[string]*LineItem{}
	items := []*LineItem{}

	for {
		if item, err := reader.Read(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if item.Type == separator {
			break
		} else if item.Type == comment {
			items = append(items, item)
		} else if item.Type == valueOnly {
			items = append(items, item)
		} else if _, contains := keys[strings.ToLower(item.Key)]; contains {
			return nil, errors.New("Malformed file--the key already exists.")
		} else {
			keys[strings.ToLower(item.Key)] = item
			items = append(items, item)
		}
	}

	return &Paragraph{keys: keys, items: items}, nil
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
