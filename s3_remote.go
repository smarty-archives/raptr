package main

// Targets S3 as a remote backend.
// Should we utilize a custom HTTP transport for things like connection pooling, keep-alive
// and custom timeouts?
type S3Remote struct {
	hostname   string // (s3-external-1.amazonaws.com by default unless region name specified via the command line)
	bucketName string
	pathPrefix string
	accessKey  string
	secretKey  string
}

func NewS3Remote(regionName, bucketName, pathPrefix, accessKey, secretKey string) *S3Remote {
	return &S3Remote{
		hostname:   resolveHostname(regionName) + ".amazonaws.com",
		bucketName: bucketName,
		pathPrefix: pathPrefix,
		accessKey:  accessKey,
		secretKey:  secretKey,
	}
}
func resolveHostname(region string) string {
	region = strings.TrimSpaces(region, " ")
	region = strings.ToLower(region)

	switch region {
	case "us-west-1":
		return "s3-us-west-1"
	case "us-west-2":
		return "s3-us-west-2"
	case "eu-west-1":
		return "s3-eu-west-1"
	case "ap-southeast-1":
		return "s3-ap-southeast-1"
	case "ap-southeast-2":
		return "s3-ap-southeast-2"
	case "ap-northeast-1":
		return "s3-ap-northeast-1"
	case "sa-east-1":
		return "s3-sa-east-1"
	default:
		return "s3-external-1"
	}

}

func (this *S3Remote) Get(request GetRequest) GetResponse {
	return GetResponse{}
}
func (this *S3Remote) Put(request PutRequest) PutResponse {
	return PutResponse{}
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
