package manifest

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestNameParser_DebianBinaryPackage(t *testing.T) { // TODO: fix broken test
	parsed := ParseFilename("/path/to/package/files/RAPTR_1.0.7-1~trusty_amd64.DEB")
	assert := assertions.New(t)
	assert.So(parsed, should.NotBeNil)
	assert.So(parsed.Name, should.Equal, "RAPTR")
	assert.So(parsed.Version, should.Equal, "1.0.7-1~trusty")
	assert.So(parsed.Architecture, should.Equal, "amd64")
	assert.So(parsed.Container, should.Equal, "deb")
}

func TestNameParser_DebianSourceCodePackage(t *testing.T) {
	parsed := ParseFilename("/path/to/package/files/raptr_1.0.7-1~trusty.DSC")
	assert := assertions.New(t)
	assert.So(parsed, should.NotBeNil)
	assert.So(parsed.Name, should.Equal, "raptr")
	assert.So(parsed.Version, should.Equal, "1.0.7-1~trusty")
	assert.So(parsed.Architecture, should.Equal, "source")
	assert.So(parsed.Container, should.Equal, "dsc")
}

func TestNameParser_ItWillIgnoreOtherFiles(t *testing.T) {
	assert := assertions.New(t)
	assert.So(ParseFilename(""), should.BeNil)
	assert.So(ParseFilename("malformed-not-enough-parts_1.0.1.deb"), should.BeNil)
	assert.So(ParseFilename("malformed_too_many_parts_1.0.1.deb"), should.BeNil)
	assert.So(ParseFilename("malformed-not-enough-parts.dsc"), should.BeNil)
	assert.So(ParseFilename("malformed_too_many_parts.dsc"), should.BeNil)
}
