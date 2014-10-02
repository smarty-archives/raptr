package remotes

// FUTURE
// Provides the ability to utilize a local/fast backend as well as a remote/slow backend.
// Any kind of concurrency-related error would invalidate the particular local item.
// Furthermore, any PUT- and DELETE-based operations should update the local/fast copy.
// Any 404/not-found operations from local/fast should forward to remote/slow.
type CachedRemote struct {
	fast Remote
	slow Remote
}

func NewCachedRemote(fast, slow Remote) *CachedRemote {
	return &CachedRemote{fast: fast, slow: slow}
}

func (this *CachedRemote) Put(request PutRequest) PutResponse {
	return this.slow.Put(request)
}
func (this *CachedRemote) Get(request GetRequest) GetResponse {
	return this.slow.Get(request)
}
func (this *CachedRemote) List(request ListRequest) ListResponse {
	return this.slow.List(request)
}
func (this *CachedRemote) Head(request HeadRequest) HeadResponse {
	return this.slow.Head(request)
}
func (this *CachedRemote) Delete(request DeleteRequest) DeleteResponse {
	return this.slow.Delete(request)
}
