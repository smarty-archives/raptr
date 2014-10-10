package manifest

import (
	"bufio"
	"errors"
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
	} else if line, err := parse(unparsed); err == pgpMessagePreamble {
		this.reader.ReadString('\n') // e.g. Hash: SHA512
		return parse("")             // separator
	} else if err == pgpSignatureBlock {
		return nil, io.EOF
	} else {
		return line, err
	}
}
func parse(unparsed string) (*LineItem, error) {
	if strings.HasSuffix(unparsed, "\r") {
		unparsed = unparsed[0 : len(unparsed)-1]
	}
	if strings.HasSuffix(unparsed, "\n") {
		unparsed = unparsed[0 : len(unparsed)-1]
	}

	if len(strings.TrimSpace(unparsed)) == 0 {
		return NewLine("", "")
	} else if strings.HasPrefix(unparsed, "#") {
		return NewLine("", unparsed)
	} else if strings.HasPrefix(unparsed, " ") || strings.HasPrefix(unparsed, "\t") {
		return NewLine("", unparsed)
	} else if colonIndex := strings.Index(unparsed, ":"); colonIndex >= 0 {
		return NewLine(unparsed[0:colonIndex], unparsed[colonIndex+1:])
	} else if unparsed == "-----BEGIN PGP SIGNED MESSAGE-----" {
		return nil, pgpMessagePreamble
	} else if unparsed == "-----BEGIN PGP SIGNATURE-----" {
		return nil, pgpSignatureBlock
	} else {
		return nil, errors.New("Malformed input")
	}
}

var pgpMessagePreamble = errors.New("PGP message preamble")
var pgpSignatureBlock = errors.New("PGPG signature block")
