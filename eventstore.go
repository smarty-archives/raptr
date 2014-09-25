package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type EventStore struct {
	sequence uint64
	handle   *os.File
	reader   *bufio.Reader
}

func NewEventStore(filename string) (*EventStore, error) {
	if handle, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		return &EventStore{
			handle: handle,
			reader: bufio.NewReaderSize(handle, bufferSize),
		}, nil
	}
}

func (this *EventStore) Next() (*EventStoreRecord, error) {
	if metaText, err := this.reader.ReadString('\n'); err != nil {
		return nil, err
	} else if meta, err := ParseMetaRecord(metaText); err != nil {
		return nil, err
	} else if payload, err := this.reader.ReadBytes('\n'); err != nil {
		return nil, err
	} else if len(payload) == 0 {
		return nil, err
	} else if err := json.Unmarshal(payload, meta.Message); err != nil {
		return nil, err
	} else if meta.Sequence != this.sequence+1 {
		return nil, errors.New("Record out of sequence")
	} else {
		this.sequence++
		return meta, nil
	}
}

func (this *EventStore) Close() error {
	return this.handle.Close()
}

const bufferSize = 1024 * 1024
