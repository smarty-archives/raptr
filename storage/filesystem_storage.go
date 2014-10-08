package storage

import (
	"bufio"
	"crypto/md5"
	"io"
	"os"
	"path"
	"path/filepath"
)

type FilesystemStorage struct{}

func NewFilesystemStorage() *FilesystemStorage {
	return &FilesystemStorage{}
}

func (this *FilesystemStorage) Put(request PutRequest) PutResponse {
	panic("Not implemented")
}
func (this *FilesystemStorage) Head(request HeadRequest) HeadResponse {
	if fullPath, err := filepath.Abs(request.Path); err != nil {
		return HeadResponse{Path: request.Path, Error: err}
	} else if info, err := os.Stat(fullPath); err != nil {
		return HeadResponse{Path: request.Path, Error: err}
	} else if md5hash, err := computeMD5(fullPath); err != nil {
		return HeadResponse{Path: request.Path, Error: err}
	} else {
		return HeadResponse{
			Path:    request.Path,
			MD5:     md5hash,
			Created: info.ModTime(),
			Length:  uint64(info.Size()),
		}
	}
}
func (this *FilesystemStorage) Get(request GetRequest) GetResponse {
	if fullPath, err := filepath.Abs(request.Path); err != nil {
		return GetResponse{Path: request.Path, Error: err}
	} else if info, err := os.Stat(fullPath); err != nil {
		return GetResponse{Path: request.Path, Error: err}
	} else if md5hash, err := computeMD5(fullPath); err != nil {
		return GetResponse{Path: request.Path, Error: err}
	} else if handle, err := os.Open(fullPath); err != nil {
		return GetResponse{Path: request.Path, Error: err}
	} else {
		return GetResponse{
			Path:     request.Path,
			Contents: handle, // implements Read/Seek/Close--application must close
			MD5:      md5hash,
			Created:  info.ModTime(),
			Length:   uint64(info.Size()),
		}
	}
}
func (this *FilesystemStorage) List(request ListRequest) ListResponse {
	if fullPath, err := filepath.Abs(request.Path); err != nil {
		return ListResponse{Path: request.Path, Error: err}
	} else {
		return ListResponse{Path: request.Path, Items: walk(fullPath)}
	}
}
func (this *FilesystemStorage) Delete(request DeleteRequest) DeleteResponse {
	panic("Not implemented")
}

func computeMD5(fullPath string) ([]byte, error) {
	handle, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	defer handle.Close()
	reader := bufio.NewReaderSize(handle, 1024*1024*16) // 16MB buffer when reading
	hasher := md5.New()

	if _, err := io.Copy(hasher, reader); err != nil {
		return nil, err
	} else {
		return hasher.Sum(nil)[:], nil
	}
}
func walk(fullPath string) []ListItem {
	// TODO: this should only walk the current directory, it shouldn't decend
	// into suboordinate directories...???
	items := []ListItem{}
	filepath.Walk(fullPath, func(directory string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		} else if err == nil {
			items = append(items, ListItem{
				Path:    path.Join(directory, info.Name()),
				Created: info.ModTime(),
				Length:  uint64(info.Size()),
			})
		}
		return nil
	})
	return items
}
