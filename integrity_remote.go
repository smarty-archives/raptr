package main

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
	return this.inner.Get(request) // verify actual MD5 vs expected MD5 as well as contents vs actual/expected
}
func (this *IntegrityRemote) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request) // verify actual MD5 vs expected MD5
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
