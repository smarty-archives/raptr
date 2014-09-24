package main

import "net/url"

type (
	RepositoryCreated struct {
		Origin string `json:"origin"`
	}

	BackendAdded struct {
		Name   string  `json:"name"`
		Remote url.URL `json:"remote"`
	}
	BackendRemoved struct {
		Name string `json:"name"`
	}

	DistributionAdded struct {
		Name string `json:"name"`
	}
	DistributionRemoved struct {
		Name string `json:"name"`
	}

	CategoryAdded struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	CategoryRemoved struct {
		Name string `json:"name"`
	}
)

type (
	PackageAdded struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		Category     string `json:"category"`
		Architecture string `json:"architecture"`
		ControlFile  string `json:"control"`
		Filename     string `json:"filename"`
		MD5          []byte `json:"md5"`
	}
	PackageRemoved struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	PackageLinked struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		Distribution string `json:"distribution"`
	}
	PackageUnlinked struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		Distribution string `json:"distribution"`
	}
)

type (
	BundleAdded struct {
		Name     string       `json:"name"`
		Version  string       `json:"version"`
		Category string       `json:"category"`
		Files    []BundleFile `json:"files"`
	}
	BundleRemoved struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	BundleLinked struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		Distribution string `json:"distribution"`
	}
	BundleUnlinked struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		Distribution string `json:"distribution"`
	}
	BundleFile struct {
		Package      string `json:"name"`
		Version      string `json:"version"`
		Architecture string `json:"architecture"`
		ControlFile  string `json:"control"`
		Filename     string `json:"filename"`
		MD5          []byte `json:"md5"`
	}
)
