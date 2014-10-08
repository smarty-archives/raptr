package manifest

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

type Checksum struct {
	Contents []byte
	MD5      []byte
	SHA1     []byte
	SHA256   []byte
	SHA512   []byte
}

func Compute(contents []byte) *Checksum {
	md5sum := md5.Sum(contents)
	sha1sum := sha1.Sum(contents)
	sha256sum := sha256.Sum256(contents)
	sha512sum := sha512.Sum512(contents)

	return &Checksum{
		MD5:    md5sum[:],
		SHA1:   sha1sum[:],
		SHA256: sha256sum[:],
		SHA512: sha512sum[:],
	}
}
