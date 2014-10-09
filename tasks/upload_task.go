package tasks

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/smartystreets/raptr/manifest"
	"github.com/smartystreets/raptr/storage"
)

type UploadTask struct {
	remote storage.Storage
	multi  *storage.MultiStorage
}

func NewUploadTask(remote storage.Storage) *UploadTask {
	return &UploadTask{remote: remote, multi: storage.NewMultiStorage(remote)}
}

// the list of files to upload along with integrity checks, etc.
// will be done during the validation phase (prior to this)
// such that any concurrency-related errors that might occur here
// won't cause the validation to be re-run
func (this *UploadTask) Upload(category, bundle, version string, packages []manifest.LocalPackage) error {
	manifestPath := manifest.BuildPath(category, bundle, version)
	log.Println("[INFO] Downloading manifest from", manifestPath)
	manifestResponse := this.remote.Get(storage.GetRequest{Path: manifestPath})
	if manifestFile, err := parseManifestResponse(manifestResponse, category, bundle, version); err != nil {
		return err // unable to access or parse remote manifest
	} else if err := this.uploadPackages(packages, manifestFile); err == noUploadedFilesError {
		log.Println("[INFO] No files uploaded, skipped uploading manifest")
		return nil
	} else if err != nil {
		return err // one or more file uploads failed
	} else {
		log.Println("[INFO] Uploading updated manifest file:", manifestPath)
		payload := manifestFile.Bytes()
		contents := storage.NewReader(payload)
		request := storage.PutRequest{Path: manifestPath, Contents: contents, Length: uint64(len(payload))}
		return this.remote.Put(request).Error
	}
}

func parseManifestResponse(response storage.GetResponse, category, bundle, version string) (*manifest.ManifestFile, error) {
	if response.Error != nil && os.IsNotExist(response.Error) {
		return manifest.NewManifestFile(category, bundle, version), nil
	} else if response.Error != nil {
		return nil, response.Error
	} else if parsed, err := manifest.ParseManifest(response.Contents, category, bundle, version); err != nil {
		return nil, err
	} else {
		log.Println("[INFO] Parsed manifest")
		return parsed, nil
	}
}
func (this *UploadTask) uploadPackages(packages []manifest.LocalPackage, manifestFile *manifest.ManifestFile) error {
	puts := []storage.PutRequest{}

	for _, pkg := range packages {
		if added, err := manifestFile.Add(pkg); err != nil {
			return err // problem adding the file to the manifest, e.g. integrity or permissions errors, etc.
		} else if !added {
			log.Printf("[INFO] The file '%s' is already contained in the manifest--SKIPPING.\n", pkg.Name())
		} else {
			for _, file := range pkg.Files() {
				targetPath := path.Join("/", manifestFile.Path(), file.Name)
				request := storage.PutRequest{Path: targetPath, Contents: file.Contents, MD5: file.Checksums.MD5, Length: file.Length}
				puts = append(puts, request)
				log.Println("Uploading local file to", request.Path)
			}
		}
	}

	for _, response := range this.multi.Put(puts...) {
		if response.Error != nil {
			return response.Error
		}
	}

	if len(puts) > 0 {
		return nil
	} else {
		return noUploadedFilesError
	}
}

var noUploadedFilesError = errors.New("No files have changed")
