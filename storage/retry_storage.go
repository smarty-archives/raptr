package storage

// Watches for connection- and mismatch-related errors and retries the
// operation up to a configured number of times.
type RetryStorage struct {
	inner      Storage
	maxRetries int
}

func NewRetryStorage(inner Storage, maxRetries int) *RetryStorage {
	return &RetryStorage{inner: inner, maxRetries: maxRetries}
}

func (this *RetryStorage) Put(request PutRequest) PutResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Put(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryStorage) Get(request GetRequest) GetResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Get(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryStorage) List(request ListRequest) ListResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.List(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryStorage) Head(request HeadRequest) HeadResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Head(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryStorage) Delete(request DeleteRequest) DeleteResponse {
	for attempt := 0; ; attempt++ {
		response := this.inner.Delete(request)
		if !this.canRetry(response.Error, attempt) {
			return response
		}
	}
}
func (this *RetryStorage) canRetry(err error, attempt int) bool {
	if attempt >= this.maxRetries {
		return false // too many attempts
	} else if err == ContentIntegrityError {
		return true // hash doesn't match actutal contents, retry
	} else if err == StorageUnavailableError {
		return true // remote system having problems, retry
	}

	return false
}
