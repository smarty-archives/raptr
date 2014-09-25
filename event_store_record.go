package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type EventStoreRecord struct {
	Sequence      uint64
	PayloadLength uint64
	Date          time.Time
	Message       interface{}
}

func ParseMetaRecord(meta string) (*EventStoreRecord, error) {
	if meta = strings.TrimSpace(meta); len(meta) < 3 {
		return nil, errors.New("Malformed meta record--record not found.")
	} else if split := strings.Split(meta[3:], ","); len(split) != expectedMetaElements {
		return nil, errors.New("Malformed meta record--missing leading comment slashes.")
	} else if sequence, err := strconv.ParseUint(split[sequenceIndex], 10, 64); err != nil {
		return nil, errors.New("Malformed meta record--bad sequence in meta record.")
	} else if instant, err := strconv.ParseUint(split[dateIndex], 10, 64); err != nil {
		return nil, errors.New("Malformed meta record--bad date in meta record.")
	} else if payloadLength, err := strconv.ParseUint(split[payloadLengthIndex], 10, 64); err != nil {
		return nil, errors.New("Malformed meta record--bad payload length in meta record.")
	} else if message, err := newMessage(split[typeIndex]); err != nil {
		return nil, err
	} else {
		return &EventStoreRecord{
			Sequence:      sequence,
			PayloadLength: payloadLength,
			Date:          time.Unix(int64(instant), 0),
			Message:       message,
		}, nil
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

const (
	expectedMetaElements = 4
	sequenceIndex        = 0
	dateIndex            = 1
	payloadLengthIndex   = 2
	typeIndex            = 3
)
