package config

import (
	"fmt"
	"strings"
)

type RepositoryLayout struct {
	Distributions []string `json:"distributions"`
	Categories    []string `json:"categories"`
	Architectures []string `json:"architectures"`
}

func (this RepositoryLayout) validate() error {
	if !isValidContents(this.Distributions) {
		return fmt.Errorf("The list of distributions has one or more missing or corrupted values.")
	} else if !isValidContents(this.Categories) {
		return fmt.Errorf("The list of categories has one or more missing or corrupted values.")
	} else if !isValidContents(this.Architectures) {
		return fmt.Errorf("The list of architectures has one or more missing or corrupted values.")
	} else {
		return nil
	}
}
func isValidContents(items []string) bool {
	for _, item := range items {
		if !isValidName(item) {
			return false
		}
	}

	return len(items) > 0
}
func isValidName(name string) bool {
	return len(strings.TrimSpace(name)) > 0
}
