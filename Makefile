#!/usr/bin/make -f

PACKAGE_NAME := raptr
PACKAGE_PATH := github.com/smartystreets/$(PACKAGE_NAME)

compile:
	go install "$(PACKAGE_PATH)"
freeze:
	glock save -n "$(PACKAGE_PATH)" > .dependencies
restore:
	cat .dependencies 2> /dev/null | glock sync -n "$(PACKAGE_PATH)"

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
	@echo "$(PACKAGE_NAME) ($(shell git describe)) unstable; urgency=low" > workspace/debian/changelog
	@echo "\n  * $(shell git rev-parse HEAD)\n" >> workspace/debian/changelog
	@echo " -- $(shell git --no-pager show -s --format="%an <%ae>")  $(shell git --no-pager show -s --format="%cD")" >> workspace/debian/changelog

dsc: clean tarball debianize changelog
	dpkg-source -b workspace

deb: dsc
	cd workspace && dpkg-buildpackage -b -us -uc

version:
	git tag -a "$(shell git describe 2>/dev/null | semver)" -m "" 2>/dev/null || true

release: clean version debianize changelog dsc
