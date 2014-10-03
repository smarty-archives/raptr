#!/usr/bin/make -f

PACKAGE_PATH := $(shell go list)
PACKAGE_NAME := $(notdir $(PACKAGE_PATH))
CLONE_DIR := $(PACKAGE_NAME)

compile: clean
	@go install
clean:
	@go clean -i
	@rm -rf "$(CLONE_DIR)"
	@rm -rf debian/files debian/$(PACKAGE_NAME)/ debian/$(PACKAGE_NAME).debhelper.log debian/$(PACKAGE_NAME).substvars
	@test -d "$(GOPATH)/debian" && (cd "$(GOPATH)" && rm -rf debian/files debian/$(PACKAGE_NAME)/ debian/$(PACKAGE_NAME).debhelper.log debian/$(PACKAGE_NAME).substvars) || echo "" > /dev/null
	@rm -f $(PACKAGE_NAME)_*
test:
	@go test -v ./...

restore: requires_tools
	@cat .dependencies 2> /dev/null | glock sync -n "$(PACKAGE_PATH)"
freeze: requires_tools
	@glock save -n "$(PACKAGE_PATH)" > .dependencies

version: requires_tools clean
	@git tag -a "$(shell git describe 2>/dev/null | semver)" -m "" 2>/dev/null || true
clone: requires_tools clean restore
	@clonetree --target="$(CLONE_DIR)" --makefile="$(PACKAGE_PATH)"
	@cp Makefile "$(CLONE_DIR)/src/$(PACKAGE_PATH)"
tarball: clean clone
	@tar -c "$(CLONE_DIR)" | gzip -n -9 > "$(PWD)/$(PACKAGE_NAME)_$(shell git describe).tar.gz"
package: clean restore version deb

dsc: requires_dpkg debianize
	@dpkg-source -b "$(CLONE_DIR)"
debianize: requires_scm clone
	@cp -r debian "$(CLONE_DIR)"
	@sed -i.bak 's/0\.0\.0/$(shell git describe)/' "$(CLONE_DIR)/debian/changelog"
	@sed -i.bak 's/none <none@none.com>/$(shell git --no-pager show -s --format="%an <%ae>")/' "$(CLONE_DIR)/debian/changelog"
	@sed -i.bak 's/Sat, 1 Jan 2000 00:00:00 +0000/$(shell git --no-pager show -s --format="%cD")/' "$(CLONE_DIR)/debian/changelog"
	@rm "$(CLONE_DIR)/debian/changelog.bak"
deb: requires_dpkg debianize
	@cd "$(CLONE_DIR)"; dpkg-buildpackage -us -uc
install:
# when the binary package is installed, this is where the artifact (app) will be installed on the target system.
	@mkdir -p "$(DESTDIR)/opt/$(PACKAGE_NAME)"
	@cp "$(GOPATH)/bin/$(PACKAGE_NAME)" "$(DESTDIR)/opt/$(PACKAGE_NAME)"

requires_tools: requires_scm
	@go get github.com/joliver/glock
	@go get github.com/smartystreets/go-packaging/semver
	@go get github.com/smartystreets/go-packaging/clonetree
requires_scm:
	@test -d .git || (echo "[ERROR] Operation only allowed on the original (SCM-controlled) instance of the source code." && exit 1)
requires_dpkg:
	@which dpkg > /dev/null || (echo "[ERROR] Debian-based dpkg-* commands are not available. Are you running Linux?" && exit 1)

%:
	@echo -n
