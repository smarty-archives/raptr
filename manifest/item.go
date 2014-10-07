package manifest

import (
	"errors"
	"strings"
)

type LineItem struct {
	itemType int
	key      string // unique per "paragraph" (case insensitive)
	value    string
}

func NewLine(key, value string) (*LineItem, error) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	if !isValidKey(key) {
		return nil, errors.New("The key of the line item contains invalid characters.")
	} else if !isValidValue(value) {
		return nil, errors.New("The value of the line item contains invalid characters.")
	} else {
		return &LineItem{
			itemType: determineType(key, value),
			key:      key,
			value:    value,
		}, nil
	}
}
func parse(unparsed string) (*LineItem, error) {
	if len(strings.TrimSpace(unparsed)) == 0 {
		return NewLine("", "")
	} else if strings.HasPrefix(unparsed, "#") {
		return NewLine("", unparsed)
	} else if strings.HasPrefix(unparsed, " ") || strings.HasPrefix(unparsed, "\t") {
		return NewLine("", unparsed)
	} else if colonIndex := strings.Index(unparsed, ":"); colonIndex >= 0 {
		return NewLine(unparsed[0:colonIndex], unparsed[colonIndex+1:])
	} else {
		return nil, errors.New("Malformed input")
	}
}

func isValidKey(text string) bool {
	for i := range text {
		if text[i] <= 32 || text[i] == 58 || text[i] >= 127 {
			return false // https://www.debian.org/doc/debian-policy/ch-controlfields.html
		}
	}

	return true
}
func isValidValue(text string) bool {
	if len(text) == 0 {
		return true
	} else if strings.Contains(text, "\n") {
		return false
	} else {
		return true
	}
}

func determineType(key, value string) int {
	if len(key) == 0 && len(value) == 0 {
		return Separator
	} else if strings.HasPrefix(value, "#") {
		return Comment
	} else if len(value) == 0 {
		return KeyOnly
	} else if len(key) == 0 {
		return ValueOnly
	} else {
		return KeyValue
	}
}

func (this LineItem) Type() int {
	return this.itemType
}
func (this *LineItem) Key() string {
	return this.key
}
func (this *LineItem) Value() string {
	return this.value
}

const (
	Separator = iota
	Comment
	KeyValue
	KeyOnly
	ValueOnly
)
