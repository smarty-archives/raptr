package main

// FUTURE
// Provides the ability to utilize a local/fast backend as well as a remote/slow backend.
// Any kind of concurrency-related error would invalidate the particular local item.
// Furthermore, any PUT- and DELETE-based operations should update the local/fast copy.
// Any 404/not-found operations from local/fast should forward to remote/slow.
type CachedBackend struct {
	fast Backend
	slow Backend
}

func NewCachedBackend(fast, slow Backend) *CachedBackend {
	return &CachedBackend{fast: fast, slow: slow}
}

func (this *CachedBackend) Put(request PutRequest) PutResponse {
	return this.slow.Put(request)
}
func (this *CachedBackend) Get(request GetRequest) GetResponse {
	return this.slow.Get(request)
}
func (this *CachedBackend) List(request ListRequest) ListResponse {
	return this.slow.List(request)
}
func (this *CachedBackend) Head(request HeadRequest) HeadResponse {
	return this.slow.Head(request)
}
func (this *CachedBackend) Delete(request DeleteRequest) DeleteResponse {
	return this.slow.Delete(request)
}
