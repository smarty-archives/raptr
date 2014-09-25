package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"
)

type EventStoreWriter struct {
	sequence uint64
	writer   io.Writer
}

func NewEventStoreWriter(sequence uint64, writer io.Writer) *EventStoreWriter {
	return &EventStoreWriter{sequence: sequence + 1, writer: writer}
}

func (this *EventStoreWriter) Write(record *EventStoreRecord) error {
	messageType := reflect.TypeOf(record.Message).Elem().Name()
	now := time.Now()

	if serialized, err := json.MarshalIndent(record.Message, "", "  "); err != nil {
		return err
	} else if fmt.Fprintf(this.writer, metaRecordFormat, this.sequence, now.Unix(), len(serialized), messageType); err != nil {
		return err
	} else if _, err := this.writer.Write(serialized); err != nil {
		return err
	} else if _, err := this.writer.Write([]byte("\n")); err != nil {
		return err
	} else {
		this.sequence++
		record.Date = now
		record.PayloadLength = uint64(len(serialized))
		return nil
	}
}

const metaRecordFormat = "// %d,%d,%d,%s\n"
