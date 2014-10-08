package tasks

import (
	"github.com/smartystreets/raptr/messages"
	"github.com/smartystreets/raptr/storage"
)

type UploadTask struct {
	local  storage.Storage
	remote storage.Storage
}

func NewUploadTask(local, remote storage.Storage) *UploadTask {
	return &UploadTask{local: local, remote: remote}
}

func (this *UploadTask) Upload(command messages.UploadCommand) error {
	// discover info about the files, e.g. deb/dsc, version info, etc.
	// check the integrity of any files listed in the dsc
	// should this be done during the validation phase of the command
	// such that it's not duplicated here? (e.g. during concurrency retry
	// attempts we shouldn't perform the above steps again).

	// download remote manifest (if any) to see packages are already there (based upon name) and skip them
	// upload the files
	// upload updated manifest (with concurrency checks)
}
