package main

import (
	"bufio"
	"encoding/json"
	"io"
)

type EventStoreReader struct {
	sequence uint64
	reader   *bufio.Reader
}

func NewEventStoreReader(reader io.Reader) *EventStoreReader {
	return &EventStoreReader{reader: bufio.NewReader(reader)}
}

func (this *EventStoreReader) Read() (*EventStoreRecord, error) {
	if metaText, err := this.reader.ReadString('\n'); err != nil {
		return nil, err
	} else if meta, err := ParseMetaRecord(this.sequence+1, metaText); err != nil {
		return nil, err
	} else if payload, err := this.reader.ReadBytes('\n'); err != nil {
		return nil, err
	} else if len(payload) == 0 {
		return nil, err
	} else if err := json.Unmarshal(payload, meta.Message); err != nil {
		return nil, err
	} else {
		this.sequence++
		return meta, nil
	}
}
