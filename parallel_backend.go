package main

// Takes a set of operations and performs them in parallel such that each
// operation happens side by side. The result slice from each operation will
// be returned in the same order, e.g. the first request correlates to the
// first response, the second to the second, and so forth. Furthermore,
// each set parallel operations will block until all desired operations have
// completed.
type ParallelBackend struct {
	backends []Backend
}

func NewParallelBackend(backends []Backend) *ParallelBackend {
	return &ParallelBackend{backends: backends}
}

func (this *ParallelBackend) Put(requests ...PutRequest) []PutResponse {
	return []PutResponse{}
}
func (this *ParallelBackend) Get(requests ...GetRequest) []GetResponse {
	return []GetResponse{}
}
func (this *ParallelBackend) Delete(requests ...DeleteRequest) []DeleteResponse {
	return []DeleteResponse{}
}
func (this *ParallelBackend) Head(requests ...HeadRequest) []HeadResponse {
	return []HeadResponse{}
}
func (this *ParallelBackend) List(requests ...ListRequest) []ListResponse {
	return []ListResponse{}
}
