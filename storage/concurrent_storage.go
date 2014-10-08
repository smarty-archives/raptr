package storage

import "bytes"

// Ensures that multiple writers (different processes or different machines)
// can be aware of each other to allow reconciliation of potentially conflicting changes.
// Requires "read your writes" consistency which S3 can provide--even in US Standard
// so long as we're targeting the s3-external-1.amazonaws.com when looking at US Standard
// buckets. When using any other region, we're fine.
// The desired behavior for concurrency-related errors is to restart the entire workflow
// such that all indexes are re-downloaded and the operation is re-attempted
type ConcurrentStorage struct {
	inner Storage
}

func NewConcurrentStorage(inner Storage) *ConcurrentStorage {
	return &ConcurrentStorage{inner: inner}
}

func (this *ConcurrentStorage) Put(request PutRequest) PutResponse {
	if err := this.ensureContents(request, CheckBeforePut); err != nil {
		return PutResponse{Path: request.Path, Error: err} // md5 differs from expected
	} else if this.alreadyExists(request) {
		return PutResponse{Path: request.Path} // no need to upload a duplicate file
	} else if response := this.inner.Put(request); response.Error != nil {
		return response // error trying to put contents
	} else if err := this.ensureContents(request, CheckAfterPut); err != nil {
		return PutResponse{Path: request.Path, Error: err} // md5 differs from expected
	} else {
		return response //success
	}
}
func (this *ConcurrentStorage) ensureContents(request PutRequest, concurrency int) error {
	if request.Concurrency&concurrency != concurrency {
		return nil // this isn't the concurrency level you're looking for
	} else if response := this.inner.Head(HeadRequest{Path: request.Path}); response.Error != nil {
		return response.Error // not found, permissions, unavailable, etc.
	} else if bytes.Compare(request.ExpectedMD5, response.MD5) != 0 {
		return ConcurrencyError
	} else {
		return nil
	}
}
func (this *ConcurrentStorage) alreadyExists(request PutRequest) bool {
	if len(request.MD5) == 0 {
		return false // no MD5 to help perform a conditional put
	} else if response := this.inner.Head(HeadRequest{Path: request.Path}); response.Error != nil {
		return false // doesn't exist or some other kind error
	} else {
		return bytes.Equal(request.MD5, response.MD5) // same = already exists
	}
}

func (this *ConcurrentStorage) Get(request GetRequest) GetResponse {
	return this.inner.Get(request)
}
func (this *ConcurrentStorage) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *ConcurrentStorage) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request)
}
func (this *ConcurrentStorage) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
