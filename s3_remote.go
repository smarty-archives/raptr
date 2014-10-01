package main

import (
	"encoding/hex"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)
import "github.com/smartystreets/go-aws-auth"

type S3Remote struct {
	hostname    string
	bucketName  string
	pathPrefix  string
	credentials awsauth.Credentials
}

func NewS3Remote(regionName, bucketName, pathPrefix, accessKey, secretKey string) *S3Remote {
	return &S3Remote{
		hostname:    resolveHostname(regionName) + ".amazonaws.com",
		bucketName:  bucketName,
		pathPrefix:  pathPrefix,
		credentials: awsauth.Credentials{AccessKeyID: accessKey, SecretAccessKey: secretKey},
	}
}

func resolveHostname(region string) string {
	region = strings.TrimSpace(region)
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

func (this *S3Remote) Head(operation HeadRequest) HeadResponse {
	request, _ := http.NewRequest("HEAD", this.composeURL(operation.Path), nil)
	if response, err := this.executeRequest(request); err != nil {
		// understand what kind of error this is--it should be transport level, not things like 404, 401, etc., etc.
		return HeadResponse{Path: operation.Path, Error: RemoteUnavailableError} // TODO: actual errors
	} else if response.StatusCode == http.StatusNotFound {
		return HeadResponse{Path: operation.Path, Error: FileNotFoundError}
	} else {
		return HeadResponse{
			Path:    operation.Path,
			MD5:     parseMD5(response.Header.Get("ETag")),
			Created: parseDate(response.Header.Get("Last-Modified")),
			Length:  parseLength(response.Header.Get("Content-Length")),
			Error:   nil,
		}
	}
}
func (this *S3Remote) Get(operation GetRequest) GetResponse {
	// create a request (construct the URL)
	// sign the request
	// issue the request with appropriate timeouts, etc.
	// for gets, ensure the content integrity is okay
	// return the response
	return GetResponse{}
}

func (this *S3Remote) Put(operation PutRequest) PutResponse {
	return PutResponse{}
}
func (this *S3Remote) List(operation ListRequest) ListResponse {
	// create a request (construct the URL)
	// sign the request
	// issue the request with appropriate timeouts, etc.
	// enumerate the results on S3
	// return the response
	return ListResponse{}
}
func (this *S3Remote) Delete(operation DeleteRequest) DeleteResponse {
	return DeleteResponse{}
}
func (this *S3Remote) composeURL(file string) string {
	return "https://" + this.hostname + path.Join("/", this.bucketName, this.pathPrefix, file)
}
func (this *S3Remote) executeRequest(request *http.Request) (*http.Response, error) {
	// TODO: connection pooling? HTTP and TCP keep alive?
	// SSL negotiation? HTTP request pipelining???
	// TCP connection timeouts, SSL handshake timeouts
	// follow redirects policy should re-sign redirects if they are for s3 resources
	awsauth.Sign(request, this.credentials)
	client := http.Client{}
	return client.Do(request)
}
func parseMD5(encoded string) []byte {
	if len(encoded) > 1 && strings.HasPrefix(encoded, `"`) {
		encoded = encoded[1 : len(encoded)-1] // strip off leading and trailing quotes
	}

	parsed, _ := hex.DecodeString(encoded)
	return parsed
}
func parseDate(date string) time.Time {
	parsed, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", date)
	return parsed
}
func parseLength(length string) uint64 {
	parsed, _ := strconv.ParseUint(length, 10, 64)
	return parsed
}
