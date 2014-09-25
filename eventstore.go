package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type EventStore struct {
	handle  *os.File
	reader  *bufio.Reader
	decoder *json.Decoder
}

func NewEventStore(filename string) (*EventStore, error) {
	if handle, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		reader := bufio.NewReaderSize(handle, bufferSize)
		return &EventStore{
			handle:  handle,
			reader:  reader,
			decoder: json.NewDecoder(reader),
		}, nil
	}
}

func (this *EventStore) Next() (interface{}, error) {
	if metaText, err := this.reader.ReadString('\n'); err != nil {
		return nil, err
	} else if meta, err := ParseMeta(metaText); err != nil {
		return nil, err
	} else if err := this.decoder.Decode(meta.Message); err != nil {
		return nil, err
	} else {
		return meta.Message, nil
	}
}

func (this *EventStore) Close() error {
	return this.handle.Close()
}

const bufferSize = 1024 * 1024
