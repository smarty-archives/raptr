package main

// Enables a locally configured file system directory to serve as a remote.
type FilesystemRemote struct{}

func NewFilesystemRemote(fast, slow Remote) *FilesystemRemote {
	return &FilesystemRemote{}
}

func (this *FilesystemRemote) Put(request PutRequest) PutResponse {
	return PutResponse{}
}
func (this *FilesystemRemote) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *FilesystemRemote) List(request ListRequest) ListResponse {
	return ListResponse{}
}
func (this *FilesystemRemote) Head(request HeadRequest) HeadResponse {
	return HeadResponse{}
}
func (this *FilesystemRemote) Delete(request DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
