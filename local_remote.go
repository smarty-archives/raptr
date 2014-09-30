package main

// Targets local filesystem directory as a remote backend.
type LocalRemote struct{}

func NewLocalRemote(fast, slow Remote) *LocalRemote {
	return &LocalRemote{}
}

func (this *LocalRemote) Put(request PutRequest) PutResponse {
	return PutResponse{}
}
func (this *LocalRemote) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *LocalRemote) List(request ListRequest) ListResponse {
	return ListResponse{}
}
func (this *LocalRemote) Head(request HeadRequest) HeadResponse {
	return HeadResponse{}
}
func (this *LocalRemote) Delete(request DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
