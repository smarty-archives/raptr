package main

type Backend interface {
	Put(PutRequest) PutResponse
	Get(GetRequest) GetResponse
	Delete(DeleteRequest) DeleteResponse
	Head(HeadRequest) HeadResponse
	List(ListRequest) ListResponse
}
