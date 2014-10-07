package manifest

import (
	"errors"
	"strings"
)

type lineItem struct {
	Type  int
	Key   string // unique per "paragraph" (case insensitive)
	Value string
}

func newLine(key, value string) (*lineItem, error) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	if !isValidKey(key) {
		return nil, errors.New("The key of the line item contains invalid characters.")
	} else if !isValidValue(value) {
		return nil, errors.New("The value of the line item contains invalid characters.")
	} else {
		return &lineItem{
			Type:  determineType(key, value),
			Key:   key,
			Value: value,
		}, nil
	}
}
func parse(unparsed string) (*lineItem, error) {
	if len(strings.TrimSpace(unparsed)) == 0 {
		return newLine("", "")
	} else if strings.HasPrefix(unparsed, "#") {
		return newLine("", unparsed)
	} else if strings.HasPrefix(unparsed, " ") || strings.HasPrefix(unparsed, "\t") {
		return newLine("", unparsed)
	} else if colonIndex := strings.Index(unparsed, ":"); colonIndex >= 0 {
		return newLine(unparsed[0:colonIndex], unparsed[colonIndex+1:])
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

const (
	separator = iota
	comment
	keyValue
	keyOnly
	valueOnly
)
