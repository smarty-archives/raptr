package manifest

import (
	"errors"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (this *Writer) Write(line *LineItem) error {
	if line == nil {
		return errors.New("Line not provided.")
	} else {
		composed := []byte(line.String() + "\n")
		_, err := this.writer.Write(composed)
		return err
	}
}
