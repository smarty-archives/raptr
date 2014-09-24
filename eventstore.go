package main

import (
	"bufio"
	"os"
)

type EventStore struct {
	handle *os.File
	reader *bufio.Reader
}

func NewEventStore(filename string) (*EventStore, error) {
	if handle, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		return &EventStore{
			handle: handle,
			reader: bufio.NewReaderSize(handle, 1024*1024),
		}, nil
	}
}

func (this *EventStore) Next() (interface{}, error) {

	return nil, nil
}

func (this *EventStore) Close() error {
	return this.handle.Close()
}
