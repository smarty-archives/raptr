package storage

import "sync"

// Takes a set of operations and performs them simultaneously. The resulting
// slice from each operation will be returned in the same order as the request
// slice. Furthermore, each set of operations will block until all desired
// operations have completed.
type MultiStorage struct {
	inner Storage
}

func NewMultiStorage(inner Storage) *MultiStorage {
	return &MultiStorage{inner: inner}
}

func (this *MultiStorage) Put(requests ...PutRequest) []PutResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]PutResponse, len(requests))

	for i, request := range requests {
		go func() {
			responses[i] = this.inner.Put(request)
			waiter.Done()
		}()
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) Get(requests ...GetRequest) []GetResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]GetResponse, len(requests))

	for i, request := range requests {
		go func() {
			responses[i] = this.inner.Get(request)
			waiter.Done()
		}()
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) Delete(requests ...DeleteRequest) []DeleteResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]DeleteResponse, len(requests))

	for i, request := range requests {
		go func() {
			responses[i] = this.inner.Delete(request)
			waiter.Done()
		}()
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) Head(requests ...HeadRequest) []HeadResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]HeadResponse, len(requests))

	for i, request := range requests {
		go func() {
			responses[i] = this.inner.Head(request)
			waiter.Done()
		}()
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) List(requests ...ListRequest) []ListResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]ListResponse, len(requests))

	for i, request := range requests {
		go func() {
			responses[i] = this.inner.List(request)
			waiter.Done()
		}()
	}

	waiter.Wait()
	return responses
}
