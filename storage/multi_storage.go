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
		go func(index int, actual PutRequest) {
			responses[index] = this.inner.Put(actual)
			waiter.Done()
		}(i, request)
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) Get(requests ...GetRequest) []GetResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]GetResponse, len(requests))

	for i, request := range requests {
		go func(index int, actual GetRequest) {
			responses[index] = this.inner.Get(actual)
			waiter.Done()
		}(i, request)
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) Delete(requests ...DeleteRequest) []DeleteResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]DeleteResponse, len(requests))

	for i, request := range requests {
		go func(index int, actual DeleteRequest) {
			responses[index] = this.inner.Delete(actual)
			waiter.Done()
		}(i, request)
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) Head(requests ...HeadRequest) []HeadResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]HeadResponse, len(requests))

	for i, request := range requests {
		go func(index int, actual HeadRequest) {
			responses[index] = this.inner.Head(actual)
			waiter.Done()
		}(i, request)
	}

	waiter.Wait()
	return responses
}
func (this *MultiStorage) List(requests ...ListRequest) []ListResponse {
	var waiter sync.WaitGroup
	waiter.Add(len(requests))
	responses := make([]ListResponse, len(requests))

	for i, request := range requests {
		go func(index int, actual ListRequest) {
			responses[index] = this.inner.List(actual)
			waiter.Done()
		}(i, request)
	}

	waiter.Wait()
	return responses
}
