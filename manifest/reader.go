package manifest

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(reader)}
}

func (this *Reader) Read() (*LineItem, error) {
	if unparsed, err := this.reader.ReadString('\n'); err != nil {
		return nil, err
	} else if strings.HasSuffix(unparsed, "\r") {
		return parse(unparsed[0 : len(unparsed)-1])
	} else {
		return parse(unparsed)
	}
}
