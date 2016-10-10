package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	region     string
	bucketName string
	pathPrefix string
}

func NewS3Storage(regionName, bucketName, pathPrefix string) *S3Storage {
	return &S3Storage{
		region:     regionName,
		bucketName: bucketName,
		pathPrefix: pathPrefix,
	}
}

func (this *S3Storage) Head(operation HeadRequest) HeadResponse {
	connection := this.buildS3Client()
	input := &s3.HeadObjectInput{
		Bucket: aws.String(this.bucketName),
		Key:    aws.String(operation.Path),
	}

	if result, err := connection.HeadObject(input); err != nil {
		return HeadResponse{Path: operation.Path, Error: err}
	} else {
		return HeadResponse{
			Path:   operation.Path,
			MD5:    parseMD5(*result.ETag),
			Length: uint64(*result.ContentLength),
		}
	}
}
func (this *S3Storage) Get(operation GetRequest) GetResponse {
	buffer := []byte{}
	connection := this.buildS3Client()
	downloader := s3manager.NewDownloaderWithClient(connection)
	object := &s3.GetObjectInput{
		Bucket: aws.String(this.bucketName),
		Key:    aws.String(operation.Path),
	}

	if read, err := downloader.Download(aws.NewWriteAtBuffer(buffer), object); err != nil {
		return GetResponse{Path: operation.Path, Error: parseError(err)}
	} else {
		return GetResponse{
			Path:     operation.Path,
			MD5:      md5.New().Sum(buffer)[:],
			Length:   uint64(read),
			Contents: NewReader(buffer),
		}
	}
}

func (this *S3Storage) Put(operation PutRequest) PutResponse {
	log.Println("[INFO] Beginning multi-part upload for", path.Base(operation.Path))

	connection := this.buildS3Client()
	uploader := s3manager.NewUploaderWithClient(connection)
	uploader.PartSize = 1024 * 1024 * 100 // 100 MB
	uploader.Concurrency = 32
	uploader.LeavePartsOnError = false

	disposition := fmt.Sprintf(contentDispositionFormat, path.Base(operation.Path))
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(this.bucketName),
		Key:                  aws.String(operation.Path),
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String(disposition),
		ServerSideEncryption: aws.String(encryption),
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
	connection := this.buildS3Client()
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(this.bucketName),
		Key:    aws.String(operation.Path),
	}

	_, err := connection.DeleteObject(input)
	return DeleteResponse{Path: operation.Path, Error: err}
}

func (this *S3Storage) buildS3Client() *s3.S3 {
	fromEnv := &credentials.EnvProvider{}
	fromEC2 := &ec2rolecreds.EC2RoleProvider{}
	chainCredentials := []credentials.Provider{fromEnv, fromEC2}
	config := &aws.Config{
		Region:      aws.String(this.region),
		Credentials: credentials.NewChainCredentials(chainCredentials),
	}
	return s3.New(session.New(config))
}

func parseMD5(encoded string) []byte {
	if len(encoded) > 1 && strings.HasPrefix(encoded, `"`) && strings.HasSuffix(encoded, `"`) {
		encoded = encoded[1 : len(encoded)-1] // strip off leading and trailing quotes
	}

	parsed, _ := hex.DecodeString(encoded)
	return parsed
}
func parseError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "NoSuchBucket", "NoSuchKey":
			return FileNotFoundError
		case "InvalidAccessKeyId":
			return AccessDeniedError
		}
	}

	return err
}

var encryption = "AES256"
var contentType = "binary/octet-stream"
var contentDispositionFormat = `attachment; filename="%s"`
