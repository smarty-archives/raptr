package main

// Enables a locally configured file system directory to serve as a backend.
type LocalBackend struct{}

func NewLocalBackend(fast, slow Backend) *LocalBackend {
	return &LocalBackend{}
}

func (this *LocalBackend) Put(request PutRequest) PutResponse {
	return PutResponse{}
}
func (this *LocalBackend) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *LocalBackend) List(request ListRequest) ListResponse {
	return ListResponse{}
}
func (this *LocalBackend) Head(request HeadRequest) HeadResponse {
	return HeadResponse{}
}
func (this *LocalBackend) Delete(request DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
