package remotes

// Watches for connection- and mismatch- errors (remote/backend unavailable) and retries the
// operation a configured number of times.
type RetryRemote struct {
	inner      Remote
	maxRetries int
}

func NewRetryRemote(inner Remote, maxRetries int) *RetryRemote {
	return &RetryRemote{inner: inner, maxRetries: maxRetries}
}

func (this *RetryRemote) Put(request PutRequest) PutResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Put(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryRemote) Get(request GetRequest) GetResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Get(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryRemote) List(request ListRequest) ListResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.List(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryRemote) Head(request HeadRequest) HeadResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Head(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryRemote) Delete(request DeleteRequest) DeleteResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Delete(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryRemote) canRetry(err error, attempt int) bool {
	if attempt >= this.maxRetries {
		return false // too many attempts
	} else if err == ContentIntegrityError {
		return true // hash doesn't match actutal contents, retry
	} else if err == RemoteUnavailableError {
		return true // remote system having problems, retry
	}

	return false
}
