package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type EventStoreRecord struct {
	Sequence uint64
	Date     time.Time
	Message  interface{}
}

func ParseMetaRecord(meta string) (*EventStoreRecord, error) {
	if meta = strings.TrimSpace(meta); len(meta) < 3 {
		return nil, errors.New("Meta record not found.")
	} else if split := strings.Split(meta[3:], ","); len(split) != expectedMetaElements {
		return nil, errors.New("Malformed meta record--missing leading comment slashes")
	} else if sequence, err := strconv.ParseUint(split[sequenceIndex], 10, 64); err != nil {
		return nil, errors.New("Malformed sequence in meta record.")
	} else if instant, err := strconv.ParseInt(split[dateIndex], 10, 64); err != nil {
		return nil, errors.New("Malformed date in meta record.")
	} else if message, err := newMessage(split[typeIndex]); err != nil {
		return nil, err
	} else {
		return &EventStoreRecord{Sequence: sequence, Date: time.Unix(instant, 0), Message: message}, nil
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

func (this *EventStoreRecord) MetaRecord() string {
	messageType := reflect.TypeOf(this.Message).Elem().Name()
	return fmt.Sprintf("// %d,%d,%s\n", this.Sequence, this.Date.Unix(), messageType)
}

const (
	expectedMetaElements = 3
	sequenceIndex        = 0
	dateIndex            = 1
	typeIndex            = 2
)
