package main

import (
	"io"
	"time"
)

type (
	GetRequest struct {
		Path        string
		ExpectedMD5 []byte
	}

	GetResponse struct {
		Path     string
		MD5      []byte
		Created  time.Time
		Length   uint64
		Contents io.Reader
		Error    error
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
		MD5         []byte
		Contents    io.Reader
		Concurrency int
		Overwrite   int
	}

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

const (
	ChaosConcurrency = iota
	CheckBeforePut
	CheckAfterPut
)

const (
	OverwriteAlways = iota
	OverwriteNever
	OverwriteIfDifferentContents
)
