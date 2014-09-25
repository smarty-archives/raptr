package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type StorageMeta struct {
	Sequence uint64
	Date     time.Time
	Message  interface{}
}

func ParseMeta(meta string) (*StorageMeta, error) {
	if meta = strings.TrimSpace(meta); len(meta) < 3 {
		return nil, errors.New("Meta record not found.")
	} else if split := strings.Split(meta[3:], "|"); len(split) != expectedMetaElements {
		return nil, errors.New("Malformed meta record.")
	} else if sequence, err := strconv.ParseUint(split[sequenceIndex], base10, uintBits); err != nil {
		return nil, errors.New("Malformed sequence in meta record.")
	} else if recorded, err := time.Parse(time.RFC3339Nano, split[dateIndex]); err != nil {
		return nil, errors.New("Malformed date in meta record.")
	} else if message, err := newMessage(split[typeIndex]); err != nil {
		return nil, err
	} else {
		return &StorageMeta{Sequence: sequence, Date: recorded, Message: message}, nil
	}
}

func newMessage(name string) (interface{}, error) {
	// TODO
	switch name {
	case "BundleAdded":
		return &BundleAdded{}, nil
	}

	return nil, errors.New("Unable to find type:" + name)
}

func (this *StorageMeta) String() string {
	return ""
}

const (
	removeLeadingComment = 3
	base10               = 10
	uintBits             = 64
	expectedMetaElements = 3
	sequenceIndex        = 0
	dateIndex            = 1
	typeIndex            = 2
)
