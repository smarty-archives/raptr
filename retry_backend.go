package main

// Watches for connection-related errors (backend unavailable) and retries the
// operation a configured number of times.
type RetryBackend struct {
	inner      Backend
	maxRetries int
}

func NewRetryBackend(inner Backend, maxRetries int) *RetryBackend {
	return &RetryBackend{inner: inner, maxRetries: maxRetries}
}

func (this *RetryBackend) Put(request PutRequest) PutResponse {
	return this.inner.Put(request)
}
func (this *RetryBackend) Get(request GetRequest) GetResponse {
	return this.inner.Get(request)
}
func (this *RetryBackend) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *RetryBackend) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request)
}
func (this *RetryBackend) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
