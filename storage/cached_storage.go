package storage

// FUTURE
// Provides the ability to utilize a local/fast backend as well as a remote/slow backend.
// Any kind of concurrency-related error would invalidate the particular local item.
// Furthermore, any PUT- and DELETE-based operations should update the local/fast copy.
// Any 404/not-found operations from local/fast should forward to remote/slow.
type CachedStorage struct {
	fast Storage
	slow Storage
}

func NewCachedStorage(fast, slow Storage) *CachedStorage {
	return &CachedStorage{fast: fast, slow: slow}
}

func (this *CachedStorage) Put(request PutRequest) PutResponse {
	return this.slow.Put(request)
}
func (this *CachedStorage) Get(request GetRequest) GetResponse {
	return this.slow.Get(request)
}
func (this *CachedStorage) List(request ListRequest) ListResponse {
	return this.slow.List(request)
}
func (this *CachedStorage) Head(request HeadRequest) HeadResponse {
	return this.slow.Head(request)
}
func (this *CachedStorage) Delete(request DeleteRequest) DeleteResponse {
	return this.slow.Delete(request)
}
