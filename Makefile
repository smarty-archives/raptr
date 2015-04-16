#!/usr/bin/make -f

SOURCE_NAME := raptr
SOURCE_VERSION := 0.1
PACKAGE_NAME := github.com/smartystreets/$(SOURCE_NAME)

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
	$(eval PREFIX := $(SOURCE_VERSION)$(shell grep "native" debian/source/format > /dev/null 2>&1 && echo "." || echo "-"))
	$(eval CURRENT := $(shell git describe 2>/dev/null))
	$(eval EXPECTED := $(PREFIX)$(shell git tag -l "$(PREFIX)*" | wc -l | xargs expr -1 +))
	$(eval INCREMENTED := $(PREFIX)$(shell git tag -l "$(PREFIX)*" | wc -l | xargs expr 0 +))
	@if [ "$(CURRENT)" != "$(EXPECTED)" ]; then git tag -a "$(INCREMENTED)" -m "" 2>/dev/null || true; fi

release: clean version debianize changelog dsc

compile:
	go install "$(PACKAGE_NAME)"
freeze:
	glock save -n "$(PACKAGE_NAME)" > .dependencies
restore:
	cat .dependencies 2> /dev/null | glock sync -n "$(PACKAGE_NAME)"
