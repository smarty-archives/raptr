package manifest

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLineItem(t *testing.T) {
	Convey("When parsing a raw Debian line item", t, func() {
		Convey("It should correctly parse simple key-value lines", func() {
			parsed, err := parse("Source: nginx")
			So(parsed, ShouldResemble, &LineItem{
				Type:  keyValue,
				Key:   "Source",
				Value: "nginx",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse simple key-value lines with whitespace", func() {
			parsed, err := parse("Source:		nginx")
			So(parsed, ShouldResemble, &LineItem{
				Type:  keyValue,
				Key:   "Source",
				Value: "nginx",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse simple key-value lines website URLs", func() {
			parsed, err := parse("URL: http://google.com/")
			So(parsed, ShouldResemble, &LineItem{
				Type:  keyValue,
				Key:   "URL",
				Value: "http://google.com/",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse key-only lines", func() {
			parsed, err := parse("Source: ")
			So(parsed, ShouldResemble, &LineItem{
				Type:  keyOnly,
				Key:   "Source",
				Value: "",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse space-prefixed value-only lines", func() {
			parsed, err := parse(" nginx deb httpd optional")
			So(parsed, ShouldResemble, &LineItem{
				Type:  valueOnly,
				Key:   "",
				Value: "nginx deb httpd optional",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse tab-prefixed value-only lines", func() {
			parsed, err := parse("\tnginx deb httpd optional")
			So(parsed, ShouldResemble, &LineItem{
				Type:  valueOnly,
				Key:   "",
				Value: "nginx deb httpd optional",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse space-prefixed, value-only lines with a colon character", func() {
			parsed, err := parse(" See more at http://google.com/")
			So(parsed, ShouldResemble, &LineItem{
				Type:  valueOnly,
				Key:   "",
				Value: "See more at http://google.com/",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse tab-prefixed, value-only lines with a colon character", func() {
			parsed, err := parse("\tSee more at http://google.com/")
			So(parsed, ShouldResemble, &LineItem{
				Type:  valueOnly,
				Key:   "",
				Value: "See more at http://google.com/",
			})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse separator lines", func() {
			parsed, err := parse("")
			So(parsed, ShouldResemble, &LineItem{Type: separator})
			So(err, ShouldBeNil)
		})
		Convey("It should correctly parse comment lines", func() {
			parsed, err := parse("#comment line")
			So(parsed, ShouldResemble, &LineItem{Type: comment, Value: "#comment line"})
			So(err, ShouldBeNil)
		})
		Convey("It should reject key values containing spaces", func() {
			parsed, err := parse("Invalid Key: value")
			So(parsed, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
