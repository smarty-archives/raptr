package main

import (
	"errors"
	"strings"
)

type EventStoreRecord struct {
	Sequence uint64
	Message  interface{}
}

func ParseMetaRecord(sequence uint64, meta string) (*EventStoreRecord, error) {
	if meta = strings.TrimSpace(meta); len(meta) < 3 {
		return nil, errors.New("Malformed meta record.")
	} else if message, err := newMessage(meta[3:]); err != nil {
		return nil, err
	} else {
		return &EventStoreRecord{Sequence: sequence, Message: message}, nil
	}
}

func newMessage(name string) (interface{}, error) {
	// TODO
	switch name {
	case "BundleAdded":
		return &BundleAdded{}, nil
	case "PackageAdded":
		return &PackageAdded{}, nil
	}

	return nil, errors.New("Unable to find type:" + name)
}
