package storage

import (
	"bytes"
	"crypto/md5"
	"io"
	"log"
	"path"
)

// Ensures the integrity of all files retrieved from the remote
// by comparing the actual MD5 (if any) vs the desired or expected MD5 (if any)
// and by comparing the Contents vs the actual or expected MD5 (whichever is populated)
type IntegrityStorage struct {
	inner Storage
}

func NewIntegrityStorage(inner Storage) *IntegrityStorage {
	return &IntegrityStorage{inner: inner}
}

func (this *IntegrityStorage) Get(request GetRequest) GetResponse {
	if response := this.inner.Get(request); response.Error != nil {
		return response
	} else if hash, matched := passedIntegrityCheck(request.ExpectedMD5, response.MD5, response.Contents); matched {
		response.MD5 = hash
		return response
	} else {
		return GetResponse{Path: request.Path, Error: ContentIntegrityError}
	}
}
func (this *IntegrityStorage) Head(request HeadRequest) HeadResponse {
	if response := this.inner.Head(request); response.Error != nil {
		return response
	} else if hash, matched := passedIntegrityCheck(request.ExpectedMD5, response.MD5, nil); matched {
		response.MD5 = hash
		return response
	} else {
		return HeadResponse{Path: request.Path, Error: ContentIntegrityError}
	}
}
func passedIntegrityCheck(expected, actual []byte, contents io.ReadSeeker) ([]byte, bool) {
	if len(expected) > 0 && len(actual) > 0 && bytes.Compare(expected, actual) != 0 {
		return []byte{}, false // expected and actual hashes don't match
	} else if !contentsMatch(expected, contents) {
		return []byte{}, false // expected (if it exists), doesn't match the contents
	} else if !contentsMatch(actual, contents) {
		return []byte{}, false // actual (if it exists), doesn't match the contents
	} else if len(expected) > 0 {
		return expected, true // expected exists and matches the contents
	} else if len(actual) > 0 {
		return actual, true // actual exists and matches the contents
	} else {
		return computeHash(contents), true // no actual or expected hash to compare against
	}
}
func contentsMatch(proposed []byte, contents io.ReadSeeker) bool {
	if contents == nil {
		return true // is this really what we want?
	} else if len(proposed) == 0 {
		return true
	} else if bytes.Compare(proposed, computeHash(contents)) == 0 {
		return true
	} else {
		return false
	}
}
func computeHash(contents io.ReadSeeker) []byte {
	// raptr never does large downloads, only small ones
	// it does however, do really large uploads, but those
	// utilize the filesystem (with seek capabilities)

	hasher := md5.New()

	if contents == nil {
		return []byte{}
	} else if _, err := contents.Seek(0, 0); err != nil {
		return []byte{} // unable to seek to beginning
	} else if _, err := io.Copy(hasher, contents); err != nil {
		return []byte{} // unable to read payload
	} else if computed := hasher.Sum(nil); len(computed) == 0 {
		return []byte{} // unable to hash
	} else if _, err := contents.Seek(0, 0); err != nil {
		return []byte{} // unable to rewind the stream again
	} else {
		return computed[:]
	}
}

func (this *IntegrityStorage) Put(request PutRequest) PutResponse {
	if len(request.MD5) == 0 {
		log.Println("[INFO] Computing hash for", path.Base(request.Path))
		request.MD5 = computeHash(request.Contents)
	} else {
		log.Println("[INFO] Hash already computed for", path.Base(request.Path))
	}

	return this.inner.Put(request)
}
func (this *IntegrityStorage) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *IntegrityStorage) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
