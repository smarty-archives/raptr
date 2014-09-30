package main

// Targets S3 as a remote backend.
// Should we utilize a custom HTTP transport for things like connection pooling, keep-alive
// and custom timeouts?
type S3Backend struct{}

func NewS3Backend(fast, slow Backend) *S3Backend {
	return &S3Backend{}
}

func (this *S3Backend) Put(request PutRequest) PutResponse {
	return PutResponse{}
}
func (this *S3Backend) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *S3Backend) List(request ListRequest) ListResponse {
	return ListResponse{}
}
func (this *S3Backend) Head(request HeadRequest) HeadResponse {
	return HeadResponse{}
}
func (this *S3Backend) Delete(request DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
