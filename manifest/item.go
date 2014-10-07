package manifest

import (
	"errors"
	"fmt"
	"strings"
)

type LineItem struct {
	Type  int
	Key   string // unique per "paragraph" (case insensitive)
	Value string
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
			Type:  determineType(key, value),
			Key:   key,
			Value: value,
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
		return separator
	} else if strings.HasPrefix(value, "#") {
		return comment
	} else if len(value) == 0 {
		return keyOnly
	} else if len(key) == 0 {
		return valueOnly
	} else {
		return keyValue
	}
}

func (this *LineItem) String() string {
	switch this.Type {
	case comment:
		return this.Value
	case keyValue:
		return fmt.Sprintf("%s: %s", this.Key, this.Value)
	case keyOnly:
		return fmt.Sprintf("%s: ", this.Key) // preserve whitespace after colon
	case valueOnly:
		return fmt.Sprintf(" %s", this.Value)
	default:
		return ""
	}
}

const (
	separator = iota
	comment
	keyValue
	keyOnly
	valueOnly
)
