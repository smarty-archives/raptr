package tasks

import (
	"bytes"
	"io"
	"os/exec"
	"path"
)

type SignatureFile struct {
	fullPath string
	payload  []byte
}

// FUTURE: this could be moved into the ReleaseFile itself, e.g. signing is a behavior on it.
func SignDistributionIndex(distribution string, releaseFile []byte) (*SignatureFile, error) {
	cmd := exec.Command("gpg", "--armor", "--yes", "--detach-sign")
	cmd.Stdin = bytes.NewBuffer(releaseFile)
	if payload, err := cmd.Output(); err != nil {
		return nil, err
	} else {
		fullPath := path.Join("/dists/", distribution, "Release.gpg")
		return &SignatureFile{fullPath: fullPath, payload: payload}, nil
	}
}
func (this *SignatureFile) Path() string          { return this.fullPath }
func (this *SignatureFile) Parse(io.Reader) error { return nil }
func (this *SignatureFile) Bytes() []byte         { return this.payload }
