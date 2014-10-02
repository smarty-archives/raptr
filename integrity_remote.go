package main

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
)

// Ensures the integrity of all files retrieved from the remote
// by comparing the actual MD5 (if any) vs the desired or expected MD5 (if any)
// and by comparing the Contents vs the actual or expected MD5 (whichever is populated)
type IntegrityRemote struct {
	inner Remote
}

func NewIntegrityRemote(inner Remote) *IntegrityRemote {
	return &IntegrityRemote{inner: inner}
}

func (this *IntegrityRemote) Get(request GetRequest) GetResponse {
	if response := this.inner.Get(request); response.Error != nil {
		return response
	} else if passedIntegrityCheck(request.ExpectedMD5, response.MD5, response.Contents) {
		return response
	} else {
		return GetResponse{Path: request.Path, Error: ContentIntegrityError}
	}
}
func (this *IntegrityRemote) Head(request HeadRequest) HeadResponse {
	if response := this.inner.Head(request); response.Error != nil {
		return response
	} else if passedIntegrityCheck(request.ExpectedMD5, response.MD5, nil) {
		return response
	} else {
		return HeadResponse{Path: request.Path, Error: ContentIntegrityError}
	}
}
func passedIntegrityCheck(expected, actual []byte, contents io.ReadSeeker) bool {
	if len(expected) > 0 && len(actual) > 0 && bytes.Compare(expected, actual) != 0 {
		return false
	} else if !contentsMatch(expected, contents) {
		return false
	} else if !contentsMatch(actual, contents) {
		return false
	} else {
		return true
	}
}
func contentsMatch(hash []byte, contents io.ReadSeeker) bool {
	if contents == nil {
		return true // no contents
	} else if len(hash) == 0 {
		return true // no hash
	} else if _, err := contents.Seek(0, 0); err != nil {
		return false // unable to seek
	} else if payload, err := ioutil.ReadAll(contents); err != nil {
		return false // unable to read payload
	} else if bytes.Compare(md5.New().Sum(payload)[:], hash) != 0 {
		return false // md5 doesn't match
	} else if _, err := contents.Seek(0, 0); err != nil {
		return false // unable to rewind the stream again
	} else {
		return true
	}
}

func (this *IntegrityRemote) Put(request PutRequest) PutResponse {
	return this.inner.Put(request)
}
func (this *IntegrityRemote) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *IntegrityRemote) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
