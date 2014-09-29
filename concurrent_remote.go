package main

import (
	"bytes"
	"errors"
)

type ConcurrentRemote struct {
	inner Remote
}

func NewConcurrentRemote(inner Remote) *ConcurrentRemote {
	return &ConcurrentRemote{inner: inner}
}

func (this *ConcurrentRemote) Put(request PutRequest) PutResponse {
	if request.Concurrency == ChaosConcurrency {
		return this.inner.Put(request)
	}

	if err := checkFile(request, CheckBeforePut); err != nil {
		return PutResponse{Path: request.Path, Error: err}
	}

	if request.Concurrency&CheckAfterPut == CheckAfterPut {
	}

	return this.inner.Put(request)
}
func (this *ConcurrentRemote) checkFile(request PutRequest, concurrency int) error {
	if request.Concurrency&concurrency != concurrency {
		return nil
	}

	if request.Concurrency&CheckBeforePut == CheckBeforePut {
		response := this.inner.Head(HeadRequest{Path: request.Path})
		if response.Error != nil {
			return PutResponse{Path: request.Path, Error: response.Error}
		} else if bytes.Compare(request.ExpectedMD5, response.MD5) {
			return ConcurrencyError
		}
	}

}

func (this *ConcurrentRemote) Delete(request DeleteRequest) DeleteResponse {
	return this.inner.Delete(request)
}
func (this *ConcurrentRemote) Get(request GetRequest) GetResponse {
	return this.inner.Get(request)
}
func (this *ConcurrentRemote) List(request ListRequest) ListResponse {
	return this.inner.List(request)
}
func (this *ConcurrentRemote) Head(request HeadRequest) HeadResponse {
	return this.inner.Head(request)
}

var ConcurrencyError = errors.New("The remote file is different from what was expected.")
