package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type EventStore struct {
	handle   *os.File
	reader   *bufio.Reader
	decoder  *json.Decoder
	sequence uint64
}

func NewEventStore(filename string) (*EventStore, error) {
	if handle, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		return NewEventStoreReader(handle), nil
	}
}
func NewEventStoreReader(handle *os.File) *EventStore {
	reader := bufio.NewReaderSize(handle, bufferSize)
	return &EventStore{
		handle:  handle,
		reader:  reader,
		decoder: json.NewDecoder(reader),
	}
}

func (this *EventStore) Next() (interface{}, error) {
	if metaText, err := this.reader.ReadString('\n'); err != nil {
		return nil, err
	} else if meta, err := ParseMeta(metaText); err != nil {
		return nil, err
	} else if err := this.decoder.Decode(meta.Message); err != nil {
		return nil, err
	} else if meta.Sequence != this.sequence+1 {
		return nil, errors.New("Record out of sequence")
	} else {
		this.sequence++
		return meta.Message, nil
	}
}

func (this *EventStore) Close() error {
	return this.handle.Close()
}

const bufferSize = 1024 * 1024
