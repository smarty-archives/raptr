package manifest

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNameParser(t *testing.T) {
	Convey("When a filename is specified", t, func() {
		Convey("It should correctly parse a debian binary package.", func() {
			parsed := ParseFilename("/path/to/package/files/RAPTR_1.0.7-1~trusty_amd64.DEB")
			So(parsed, ShouldNotBeNil)
			So(parsed.Name, ShouldEqual, "raptr")
			So(parsed.Version, ShouldEqual, "1.0.7-1~trusty")
			So(parsed.Architecture, ShouldEqual, "amd64")
			So(parsed.Container, ShouldEqual, "deb")
		})
		Convey("It should correctly parse a debian source code package.", func() {
			parsed := ParseFilename("/path/to/package/files/raptr_1.0.7-1~trusty.DSC")
			So(parsed, ShouldNotBeNil)
			So(parsed.Name, ShouldEqual, "raptr")
			So(parsed.Version, ShouldEqual, "1.0.7-1~trusty")
			So(parsed.Architecture, ShouldEqual, "source")
			So(parsed.Container, ShouldEqual, "dsc")
		})
		Convey("It should not interpret other files.", func() {
			So(ParseFilename(""), ShouldBeNil)
			So(ParseFilename("malformed-not-enough-parts_1.0.1.deb"), ShouldBeNil)
			So(ParseFilename("malformed_too_many_parts_1.0.1.deb"), ShouldBeNil)
			So(ParseFilename("malformed-not-enough-parts.dsc"), ShouldBeNil)
			So(ParseFilename("malformed_too_many_parts.dsc"), ShouldBeNil)
		})
	})
}
