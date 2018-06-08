package tasks

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

type SignatureFile struct {
	fullPath string
	payload  []byte
}

func SignDistributionIndex(distribution string, releaseFile []byte) (*SignatureFile, error) {
	cmd := createCommand()
	cmd.Stdin = bytes.NewBuffer(releaseFile)
	errorBuffer := bytes.NewBuffer([]byte{})
	cmd.Stderr = errorBuffer
	if payload, err := cmd.Output(); err != nil {
		log.Printf("[ERROR] Unable to sign distribution index: %s\n", string(errorBuffer.Bytes()))
		return nil, err
	} else {
		fullPath := path.Join("/dists/", distribution, "InRelease")
		return &SignatureFile{fullPath: fullPath, payload: payload}, nil
	}
}
func createCommand() *exec.Cmd {
	passphrase := os.Getenv("GPG_PASSPHRASE")
	if len(passphrase) == 0 {
		return exec.Command("gpg", "--clearsign", "--detach-sign")
	} else {
		return exec.Command("gpg", "--clearsign", "--pinentry-mode", "--no-tty", "--batch", "--passphrase", passphrase)
	}
}

func (this *SignatureFile) Path() string          { return this.fullPath }
func (this *SignatureFile) Parse(io.Reader) error { return nil }
func (this *SignatureFile) Bytes() []byte         { return this.payload }
