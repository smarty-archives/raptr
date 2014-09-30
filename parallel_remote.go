package main

// Takes a set of operations and performs them in parallel such that each
// operation happens side by side. The result slice from each operation will
// be returned in the same order, e.g. the first request correlates to the
// first response, the second to the second, and so forth. Furthermore,
// each set parallel operations will block until all desired operations have
// completed.
type ParallelRemote struct {
	remotes []Remote
}

func NewParallelRemote(remotes []Remote) *ParallelRemote {
	return &ParallelRemote{remotes: remotes}
}

func (this *ParallelRemote) Put(requests ...PutRequest) []PutResponse {
	return []PutResponse{}
}
func (this *ParallelRemote) Get(requests ...GetRequest) []GetResponse {
	return []GetResponse{}
}
func (this *ParallelRemote) Delete(requests ...DeleteRequest) []DeleteResponse {
	return []DeleteResponse{}
}
func (this *ParallelRemote) Head(requests ...HeadRequest) []HeadResponse {
	return []HeadResponse{}
}
func (this *ParallelRemote) List(requests ...ListRequest) []ListResponse {
	return []ListResponse{}
}
