package main

// Watches for connection-related errors (backend unavailable) and retries the
// operation a configured number of times.
type RetryRemote struct {
	inner      Remote
	maxRetries int
}

func NewRetryRemote(inner Remote, maxRetries int) *RetryRemote {
	return &RetryRemote{inner: inner, maxRetries: maxRetries}
}

func (this *RetryRemote) Put(request PutRequest) PutResponse {
	return this.inner.Put(request)
}
func (this *RetryRemote) Get(request GetRequest) GetResponse {
	return this.inner.Get(request)
}
func (this *RetryRemote) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *RetryRemote) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request)
}
func (this *RetryRemote) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
