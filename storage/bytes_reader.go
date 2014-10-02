package storage

import "bytes"

type BytesReader struct {
	reader *bytes.Reader
}

func NewReader(payload []byte) *BytesReader {
	return &BytesReader{reader: bytes.NewReader(payload)}
}
func (this *BytesReader) Read(buffer []byte) (int, error) {
	return this.reader.Read(buffer)
}
func (this *BytesReader) Seek(offset int64, whence int) (int64, error) {
	return this.reader.Seek(offset, whence)
}
func (this *BytesReader) Close() error {
	return nil
}
