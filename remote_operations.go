package main

import (
	"errors"
	"io"
	"os"
	"time"
)

type (
	GetRequest struct {
		Path string
		// ExpectedMD5 []byte // empty if we don't care
	}

	// errors resulting from download:
	// 1. nil = success
	// 2. MD5 hash/content mismatch
	// 3. 404 not found / file doesn't exist
	// 4. permissions
	// 5. remote/backend unavailable
	GetResponse struct {
		Path     string
		MD5      []byte
		Created  time.Time
		Length   uint64
		Contents io.ReadSeeker // we need to be able to read the entire stream multiple times
		Error    error         // contains not found errors, backend unavailable, etc.
	}
)

type (
	HeadRequest struct {
		Path string
	}

	HeadResponse struct {
		Path    string
		MD5     []byte
		Created time.Time
		Length  uint64
		Error   error
	}
)

type (
	PutRequest struct {
		Path        string
		MD5         []byte        // empty if we don't care
		ExpectedMD5 []byte        // empty if we don't care
		Contents    io.ReadSeeker // we need to be able to read the entire stream multiple times
		Length      uint64        // for streaming large file from filesystem; []byte can be wrapped in a buffer
		Concurrency int
		Overwrite   int
	}

	// errors resulting from upload:
	// 1. nil = success
	// 2. MD5 hash/content mismatch
	// 3. concurrency mismatch (file has changed either before or after writing depending upon desired PUT concurrency)
	// 4. permissions
	// 5. remote/backend unavailable
	PutResponse struct {
		Path  string
		Error error
	}
)

type (
	DeleteRequest struct {
		Path string
	}

	DeleteResponse struct {
		Path  string
		Error error
	}
)

type (
	ListRequest struct {
		Path string
	}

	ListResponse struct {
		Path  string
		Items []ListItem
		Error error
	}
	ListItem struct {
		Path    string
		Created time.Time
		Length  uint64
		MD5     []byte
	}
)

var (
	ConcurrencyError       = errors.New("The remote file is different from what was expected.")
	ContentIntegrityError  = errors.New("The contents of the file do not match the expected hash.")
	RemoteUnavailableError = errors.New("The remote system is unavailable or not responding.")
	AccessDeniedError      = os.ErrPermission
	FileNotFoundError      = os.ErrNotExist
)

const (
	ChaosConcurrency = 0
	CheckBeforePut   = 1 << iota
	CheckAfterPut    = 1 << iota
)

const (
	OverwriteAlways = iota
	OverwriteNever
	OverwriteIfDifferentContents
)
