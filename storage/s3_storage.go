package storage

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/smartystreets/go-aws-auth"
)

type S3Storage struct {
	region     string
	hostname   string
	bucketName string
	pathPrefix string
	client     *http.Client
}

func NewS3Storage(regionName, bucketName, pathPrefix string) *S3Storage {
	return &S3Storage{
		region:     regionName,
		hostname:   resolveHostname(regionName) + ".amazonaws.com",
		bucketName: bucketName,
		pathPrefix: pathPrefix,
		client:     buildClient(),
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
func buildClient() *http.Client {
	// TODO:
	// 1. TCP connection timeouts
	// 2. SSL handshake timeouts
	// 3. HTTP request pipelining
	// 4. HTTP keep alive
	// 5. TCP keep alive
	// 6. connection pooling
	// 7. redirect policy--don't follow any?
	return &http.Client{}
}

func (this *S3Storage) Head(operation HeadRequest) HeadResponse {
	request := this.newRequest("HEAD", operation.Path, nil)
	if response, err := this.executeRequest(request); err != nil {
		return HeadResponse{Path: operation.Path, Error: err}
	} else {
		header := response.Header
		return HeadResponse{
			Path:    operation.Path,
			MD5:     parseMD5(header.Get("ETag")),
			Created: parseDate(header.Get("Last-Modified")),
			Length:  parseLength(header.Get("Content-Length")),
		}
	}
}
func (this *S3Storage) Get(operation GetRequest) GetResponse {
	request := this.newRequest("GET", operation.Path, nil)
	if response, err := this.executeRequest(request); err != nil {
		return GetResponse{Path: operation.Path, Error: err}
	} else if payload, err := ioutil.ReadAll(response.Body); err != nil {
		// io.ReadAll isn't a problem here because we're doing very small downloads
		// e.g. indexes and signature files
		response.Body.Close()
		return GetResponse{Path: operation.Path, Error: StorageUnavailableError}
	} else {
		header := response.Header
		return GetResponse{
			Path:     operation.Path,
			MD5:      parseMD5(header.Get("ETag")),
			Created:  parseDate(header.Get("Last-Modified")),
			Length:   parseLength(header.Get("Content-Length")),
			Contents: NewReader(payload),
		}
	}
}

func (this *S3Storage) Put(operation PutRequest) PutResponse {
	if operation.Length < multiPartSize {
		return this.putSingle(operation)
	} else {
		return this.putMulti(operation)
	}
}
func (this *S3Storage) putSingle(operation PutRequest) PutResponse {
	request := this.newRequest("PUT", operation.Path, ioutil.NopCloser(operation.Contents))
	request.ContentLength = int64(operation.Length) // TODO: when this is zero (empty files) the request uses "Transfer-Encoding: Chunked"?!
	request.Header.Set("Expect", "100-continue")    // send headers before body
	request.Header.Set("x-amz-server-side-encryption", encryption)
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Content-Disposition", fmt.Sprintf(contentDispositionFormat, path.Base(operation.Path))) // correct names for wget/curl/etc
	if len(operation.MD5) > 0 {
		request.Header.Set("Content-Md5", base64.StdEncoding.EncodeToString(operation.MD5))
	}
	_, err := this.executeRequest(request)
	return PutResponse{Path: operation.Path, MD5: operation.MD5, Error: err}
}
func (this *S3Storage) putMulti(operation PutRequest) PutResponse {
	log.Println("[INFO] Beginning multi-part upload for", path.Base(operation.Path))
	manager := s3manager.NewUploader(&s3manager.UploadOptions{
		PartSize:          1024 * 1024 * 100,
		Concurrency:       32,
		LeavePartsOnError: false,
		S3: s3.New(&aws.Config{
			Credentials: credentials.NewChainCredentials([]credentials.Provider{&credentials.EnvProvider{}, &credentials.EC2RoleProvider{}}),
			Region:      this.region,
		}),
	})

	disposition := fmt.Sprintf(contentDispositionFormat, path.Base(operation.Path))
	_, err := manager.Upload(&s3manager.UploadInput{
		Bucket:               &this.bucketName,
		Key:                  &operation.Path,
		ContentType:          &contentType,
		ContentDisposition:   &disposition,
		ServerSideEncryption: &encryption,
		Body:                 operation.Contents,
	})

	if err != nil {
		log.Printf("[ERROR] Multi-part PUT Request Error for [%s]: [%s]\n", path.Base(operation.Path), err)
		return PutResponse{Path: operation.Path, MD5: operation.MD5, Error: err}
	} else {
		// nil MD5 skips concurrency check (e.g. multiple writers)
		// but there's no obvious way to get it
		return PutResponse{Path: operation.Path, MD5: []byte{}, Error: nil}
	}
}

func (this *S3Storage) List(operation ListRequest) ListResponse {
	// this will be at least one request until we've gathered everything locally
	panic("Not implemented")
}
func (this *S3Storage) Delete(operation DeleteRequest) DeleteResponse {
	request := this.newRequest("DELETE", operation.Path, nil)
	_, err := this.executeRequest(request)
	return DeleteResponse{Path: operation.Path, Error: err}
}

func (this *S3Storage) newRequest(method, requestPath string, body io.Reader) *http.Request {
	url := "https://" + this.hostname + path.Join("/", this.bucketName, this.pathPrefix, requestPath)
	request, _ := http.NewRequest(method, url, body)
	return request
}
func (this *S3Storage) executeRequest(request *http.Request) (*http.Response, error) {
	// don't use Sign4, it reads and replaces the body which can results in out of memory problems
	// and it can also affect retry
	awsauth.SignS3(request)
	if response, err := this.client.Do(request); err != nil {
		log.Printf("[ERROR] %s Request Error for [%s]: [%s]\n", request.Method, path.Base(request.URL.Path), err)
		return nil, StorageUnavailableError
	} else {
		return response, parseError(response.StatusCode, response.Body)
	}
}
func parseMD5(encoded string) []byte {
	if len(encoded) > 1 && strings.HasPrefix(encoded, `"`) && strings.HasSuffix(encoded, `"`) {
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
func parseError(statusCode int, body io.ReadCloser) error {
	if statusCode == http.StatusOK {
		return nil
	} else if statusCode == http.StatusBadRequest { // 400
		log.Println("[WARN] Bad HTTP Request:", parseBody(body))
		return ContentIntegrityError
	} else if statusCode == http.StatusUnauthorized { // 401
		return AccessDeniedError
	} else if statusCode == http.StatusForbidden { // 403
		return AccessDeniedError
	} else if statusCode == http.StatusNotFound { // 404
		return FileNotFoundError
	} else {
		log.Println("[WARN] Failed HTTP Request:", parseBody(body))
		return StorageUnavailableError
	}
}
func parseBody(body io.ReadCloser) string {
	raw, _ := ioutil.ReadAll(body)
	return string(raw)
}

const (
	multiPartSize = 5 * 1024 * 1024 * 1024
)

var encryption = "AES256"
var contentType = "binary/octet-stream"
var contentDispositionFormat = `attachment; filename="%s"`
