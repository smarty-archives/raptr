package main

import "bytes"

// Ensures that multiple writers (different processes or different machines)
// can be aware of each other to allow reconciliation of potentially conflicting changes.
// Requires "read your writes" consistency which S3 can provide--even in US Standard
// so long as we're targeting the s3-external-1.amazonaws.com when looking at US Standard
// buckets. When using any other region, we're fine.
// The desired behavior for concurrency-related errors is to restart the entire workflow
// such that all indexes are re-downloaded and the operation is re-attempted
type ConcurrentRemote struct {
	inner Remote
}

func NewConcurrentRemote(inner Remote) *ConcurrentRemote {
	return &ConcurrentRemote{inner: inner}
}

func (this *ConcurrentRemote) Put(request PutRequest) PutResponse {
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
func (this *ConcurrentRemote) ensureContents(request PutRequest, concurrency int) error {
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

func (this *ConcurrentRemote) Get(request GetRequest) GetResponse {
	return this.inner.Get(request)
}
func (this *ConcurrentRemote) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *ConcurrentRemote) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request)
}
func (this *ConcurrentRemote) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
