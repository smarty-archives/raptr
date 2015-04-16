#!/usr/bin/make -f

SOURCE_NAME := raptr
SOURCE_VERSION := 0.1
PACKAGE_NAME := github.com/smartystreets/$(SOURCE_NAME)

compile:
	go install "$(PACKAGE_NAME)"
freeze:
	glock save -n "$(PACKAGE_NAME)" > .dependencies
restore:
	cat .dependencies 2> /dev/null | glock sync -n "$(PACKAGE_NAME)"

clean:
	rm -rf workspace *.tar.?z *.dsc *.deb *.changes

prepare: clean restore
	mkdir -p workspace
	cp Releasefile workspace/Makefile
	clonetree --target=workspace

tarball: prepare

debianize:
	mkdir -p workspace
	cp -r debian workspace

changelog: debianize
	@echo "$(SOURCE_NAME) ($(shell git describe)) unstable; urgency=low" > workspace/debian/changelog
	@echo "\n  * $(shell git rev-parse HEAD)\n" >> workspace/debian/changelog
	@echo " -- $(shell git --no-pager show -s --format="%an <%ae>")  $(shell git --no-pager show -s --format="%cD")" >> workspace/debian/changelog

dsc: clean tarball debianize changelog
	dpkg-source -b workspace

deb: dsc
	cd workspace && dpkg-buildpackage -b -us -uc

version:
	git tag -a "$(shell git describe 2>/dev/null | semver)" -m "" 2>/dev/null || true

release: clean version debianize changelog dsc
