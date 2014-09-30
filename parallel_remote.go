package main

import "sync"

// Takes a set of operations and performs them in parallel such that each
// operation happens side by side. The result slice from each operation will
// be returned in the same order, e.g. the first request correlates to the
// first response, the second to the second, and so forth. Furthermore,
// each set parallel operations will block until all desired operations have
// completed.
type ParallelRemote struct {
	inner Remote
}

func NewParallelRemote(inner Remote) *ParallelRemote {
	return &ParallelRemote{inner: inner}
}

func (this *ParallelRemote) Put(requests ...PutRequest) []PutResponse {
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
func (this *ParallelRemote) Get(requests ...GetRequest) []GetResponse {
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
func (this *ParallelRemote) Delete(requests ...DeleteRequest) []DeleteResponse {
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
func (this *ParallelRemote) Head(requests ...HeadRequest) []HeadResponse {
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
func (this *ParallelRemote) List(requests ...ListRequest) []ListResponse {
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
