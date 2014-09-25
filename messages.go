package main

import "net/url"

type (
	RepositoryCreated struct {
		Origin string `json:"origin,omitempty"`
	}

	BackendAdded struct {
		Name   string  `json:"name,omitempty"`
		Remote url.URL `json:"remote,omitempty"`
	}
	BackendRemoved struct {
		Name string `json:"name,omitempty"`
	}

	DistributionAdded struct {
		Name string `json:"name,omitempty"`
	}
	DistributionRemoved struct {
		Name string `json:"name,omitempty"`
	}

	CategoryAdded struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}
	CategoryRemoved struct {
		Name string `json:"name,omitempty"`
	}
)

type (
	PackageAdded struct {
		Name         string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		Category     string `json:"category,omitempty"`
		Architecture string `json:"architecture,omitempty"`
		ControlFile  string `json:"control,omitempty"`
		Filename     string `json:"filename,omitempty"`
		MD5          []byte `json:"md5,omitempty"`
	}
	PackageRemoved struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
	}
	PackageLinked struct {
		Name         string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		Distribution string `json:"distribution,omitempty"`
	}
	PackageUnlinked struct {
		Name         string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		Distribution string `json:"distribution,omitempty"`
	}
)

type (
	BundleAdded struct {
		Name     string       `json:"name,omitempty"`
		Version  string       `json:"version,omitempty"`
		Category string       `json:"category,omitempty"`
		Files    []BundleFile `json:"files,omitempty"`
	}
	BundleRemoved struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
	}
	BundleLinked struct {
		Name         string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		Distribution string `json:"distribution,omitempty"`
	}
	BundleUnlinked struct {
		Name         string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		Distribution string `json:"distribution,omitempty"`
	}
	BundleFile struct {
		Package      string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		Architecture string `json:"architecture,omitempty"`
		ControlFile  string `json:"control,omitempty"`
		Filename     string `json:"filename,omitempty"`
		MD5          []byte `json:"md5,omitempty"`
	}
)
