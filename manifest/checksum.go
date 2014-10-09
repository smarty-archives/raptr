package manifest

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"io"
)

type Checksum struct {
	MD5    []byte
	SHA1   []byte
	SHA256 []byte
	SHA512 []byte
}

func ComputeChecksums(reader io.Reader) (Checksum, error) {
	md5hash := md5.New()
	sha1hash := sha1.New()
	sha256hash := sha256.New()
	sha512hash := sha512.New()

	writer := io.MultiWriter(md5hash, sha1hash, sha256hash, sha512hash)
	if _, err := io.Copy(writer, reader); err != nil {
		return Checksum{}, err
	} else {
		md5sum := md5hash.Sum(nil)
		sha1sum := sha1hash.Sum(nil)
		sha256sum := sha256hash.Sum(nil)
		sha512sum := sha512hash.Sum(nil)
		return Checksum{
			MD5:    md5sum[:],
			SHA1:   sha1sum[:],
			SHA256: sha256sum[:],
			SHA512: sha512sum[:],
		}, nil
	}
}
