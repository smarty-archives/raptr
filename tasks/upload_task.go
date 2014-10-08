package tasks

import (
	"log"
	"os"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type UploadTask struct {
	local  storage.Storage
	remote storage.Storage
	multi  *storage.MultiStorage
}

func NewUploadTask(local, remote storage.Storage) *UploadTask {
	return &UploadTask{local: local, remote: remote, multi: storage.NewMultiStorage(remote)}
}

// the list of files to upload along with integrity checks, etc.
// will be done during the validation phase (prior to this)
// such that any concurrency-related errors that might occur here
// won't cause the validation to be re-run
func (this *UploadTask) Upload(category, bundle, version string, files []string) error {
	manifestPath := manifest.BuildPath(category, bundle, version)
	manifestResponse := this.remote.Get(storage.GetRequest{Path: manifestPath})
	if manifestFile, err := parseManifestResponse(manifestResponse); err != nil {
		return err // unable to access remote manifest
	} else if err := this.uploadBundleFiles(files, manifestFile); err != nil {
		return err // one or more file uploads failed
	} else {
		contents := storage.NewReader(manifestFile.Bytes())
		request := storage.PutRequest{Path: manifestPath, Contents: contents}
		return this.remote.Put(request).Error
	}
}

func parseManifestResponse(response storage.GetResponse) (*manifest.ManifestFile, error) {
	if response.Error != nil && os.IsNotExist(response.Error) {
		return manifest.NewManifestFile(), nil
	} else if response.Error != nil {
		return nil, response.Error
	} else if parsed, err := manifest.ParseManifest(response.Contents); err != nil {
		return nil, err
	} else {
		return parsed, nil
	}
}
func (this *UploadTask) uploadBundleFiles(files []string, manifestFile *manifest.ManifestFile) error {
	uploads := []storage.PutRequest{}

	for _, file := range files {
		if !manifestFile.Add(file) {
			log.Println()
		} else {
			uploads = append(uploads, storage.PutRequest{}) // TODO
		}
	}

	for _, response := range this.multi.Put(uploads...) {
		if response.Error != nil {
			return response.Error
		}
	}

	return nil
}
