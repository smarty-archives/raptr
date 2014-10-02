package storage

// Enables a locally configured file system directory to serve as a remote.
type FilesystemStorage struct {
	filenameAppendMD5 bool
}

func NewFilesystemStorage(fast, slow Storage) *FilesystemStorage {
	return &FilesystemStorage{}
}

func (this *FilesystemStorage) Put(request PutRequest) PutResponse {
	return PutResponse{}
}
func (this *FilesystemStorage) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *FilesystemStorage) List(request ListRequest) ListResponse {
	return ListResponse{}
}
func (this *FilesystemStorage) Head(request HeadRequest) HeadResponse {
	return HeadResponse{}
}
func (this *FilesystemStorage) Delete(request DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
