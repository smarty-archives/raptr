package main

import (
	"net/url"
	"time"
)

type (
	RepositoryCreated struct {
		Date time.Time `json:"date"`
	}

	BackendAdded struct {
		Name   string    `json:"name"`
		Remote url.URL   `json:"remote"`
		Date   time.Time `json:"date"`
	}
	BackendRemoved struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}

	DistributionAdded struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}
	DistributionRemoved struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}

	CategoryAdded struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}
	CategoryRemoved struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}
)

type (
	PackageAdded struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Category     string    `json:"category"`
		Architecture string    `json:"architecture"`
		Filename     string    `json:"filename"`
		MD5          []byte    `json:"md5"`
		Date         time.Time `json:"date"`
	}
	PackageRemoved struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Date    time.Time `json:"date"`
	}
	PackageLinked struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Distribution string    `json:"distribution"`
		Date         time.Time `json:"date"`
	}
	PackageUnlinked struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Distribution string    `json:"distribution"`
		Date         time.Time `json:"date"`
	}
)

type (
	BundleAdded struct {
		Name     string       `json:"name"`
		Version  string       `json:"version"`
		Category string       `json:"category"`
		Files    []BundleFile `json:"files"`
		Date     time.Time    `json:"date"`
	}
	BundleRemoved struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Date    time.Time `json:"date"`
	}
	BundleLinked struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Distribution string    `json:"distribution"`
		Date         time.Time `json:"date"`
	}
	BundleUnlinked struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Distribution string    `json:"distribution"`
		Date         time.Time `json:"date"`
	}
	BundleFile struct {
		Package      string
		Version      string
		Architecture string
		Filename     string
		MD5          []byte
	}
)
