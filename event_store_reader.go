package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
	} else if meta, err := ParseMetaRecord(metaText); err != nil {
		return nil, err
	} else if meta.Sequence != this.sequence+1 {
		return nil, errors.New("Malformed meta record--incorrect sequence number.")
	} else if err := this.deserializeMessage(meta); err != nil {
		return nil, err
	} else {
		this.sequence++
		return meta, nil
	}
}
func (this *EventStoreReader) deserializeMessage(record *EventStoreRecord) error {
	limitReader := io.LimitReader(this.reader, int64(record.PayloadLength+1))
	decoder := json.NewDecoder(limitReader)
	return decoder.Decode(record.Message)
}
