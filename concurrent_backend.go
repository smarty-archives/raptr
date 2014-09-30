package main

import (
	"bytes"
	"errors"
)

// Ensures that multiple writers (different processes or different machines)
// can be aware of each other to allow reconciliation of potentially conflicting changes.
// Requires "read your writes" consistency which S3 can provide--even in US Standard
// so long as we're targeting the s3-external-1.amazonaws.com when looking at US Standard
// buckets. When using any other region, we're fine.
// The desired behavior for concurrency-related errors is to restart the entire workflow
// such that all indexes are re-downloaded and the operation is re-attempted
type ConcurrentBackend struct {
	inner Backend
}

func NewConcurrentBackend(inner Backend) *ConcurrentBackend {
	return &ConcurrentBackend{inner: inner}
}

func (this *ConcurrentBackend) Put(request PutRequest) PutResponse {
	if err := this.ensureContents(request, CheckBeforePut); err != nil {
		return PutResponse{Path: request.Path, Error: err}
	} else if response := this.inner.Put(request); response.Error != nil {
		return response
	} else if err := this.ensureContents(request, CheckAfterPut); err != nil {
		return PutResponse{Path: request.Path, Error: err}
	} else {
		return response
	}
}
func (this *ConcurrentBackend) ensureContents(request PutRequest, concurrency int) error {
	if request.Concurrency&concurrency != concurrency {
		return nil
	} else if response := this.inner.Head(HeadRequest{Path: request.Path}); response.Error != nil {
		return response.Error
	} else if bytes.Compare(request.ExpectedMD5, response.MD5) != 0 {
		return ConcurrencyError
	} else {
		return nil
	}
}

func (this *ConcurrentBackend) Get(request GetRequest) GetResponse {
	return this.inner.Get(request)
}
func (this *ConcurrentBackend) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *ConcurrentBackend) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request)
}
func (this *ConcurrentBackend) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}

var ConcurrencyError = errors.New("The remote file is different from what was expected.")
