package main

// Targets S3 as a remote backend.
// Should we utilize a custom HTTP transport for things like connection pooling, keep-alive
// and custom timeouts?
type S3Remote struct{}

func NewS3Remote(fast, slow Remote) *S3Remote {
	return &S3Remote{}
}

func (this *S3Remote) Put(request PutRequest) PutResponse {
	return PutResponse{}
}
func (this *S3Remote) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *S3Remote) List(request ListRequest) ListResponse {
	return ListResponse{}
}
func (this *S3Remote) Head(request HeadRequest) HeadResponse {
	return HeadResponse{}
}
func (this *S3Remote) Delete(request DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
