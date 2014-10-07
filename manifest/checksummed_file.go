package manifest

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

type ChecksummedFile struct {
	Contents []byte
	MD5      []byte
	SHA1     []byte
	SHA256   []byte
	SHA512   []byte
}

func NewFile(contents []byte) *ChecksummedFile {
	return &ChecksummedFile{
		Contents: contents,
		MD5:      md5.Sum(contents)[:],
		SHA1:     sha1.Sum(contents)[:],
		SHA256:   sha256.Sum256(contents)[:],
		SHA512:   sha512.Sum512(contents)[:],
	}
}
