package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type EventStoreWriter struct {
	writer  io.Writer
	encoder *json.Encoder
}

func NewEventStoreWriter(writer io.Writer) *EventStoreWriter {
	return &EventStoreWriter{writer: writer, encoder: json.NewEncoder(writer)}
}

func (this *EventStoreWriter) Write(record *EventStoreRecord) error {
	messageType := reflect.TypeOf(record.Message).Elem().Name()
	if _, err := fmt.Fprintf(this.writer, metaRecordFormat, messageType); err != nil {
		return err
	} else {
		return this.encoder.Encode(record.Message)
	}
}

const metaRecordFormat = "// %s\n"
