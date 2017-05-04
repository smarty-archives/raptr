package manifest

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestLineItemFixture(t *testing.T) {
	gunit.Run(new(LineItemFixture), t)
}

type LineItemFixture struct {
	*gunit.Fixture
}

func (this *LineItemFixture) assertParseSuccess(input string, expected *LineItem) {
	parsed, err := parse(input)
	this.So(err, should.BeNil)
	this.So(parsed, should.Resemble, expected)
}

func (this *LineItemFixture) TestSimpleKeyValueLines() {

	this.assertParseSuccess("Source: nginx", // simple key-value
		&LineItem{Type: keyValue, Key: "Source", Value: "nginx"})

	this.assertParseSuccess("Source:		nginx", // w/ whitespace
		&LineItem{Type: keyValue, Key: "Source", Value: "nginx"})

	this.assertParseSuccess("URL: http://google.com/", // URL
		&LineItem{Type: keyValue, Key: "URL", Value: "http://google.com/"})

	this.assertParseSuccess("Source: ", // Key only
		&LineItem{Type: keyOnly, Key: "Source", Value: ""})

	this.assertParseSuccess(" nginx deb httpd optional", // space prefixed value-only line
		&LineItem{Type: valueOnly, Key: "", Value: "nginx deb httpd optional"})

	this.assertParseSuccess("\tnginx deb httpd optional", // tab-prefixed value-only line
		&LineItem{Type: valueOnly, Key: "", Value: "nginx deb httpd optional"})

	this.assertParseSuccess(" See more at http://google.com/", // space prefixed value-only line w/ colon
		&LineItem{Type: valueOnly, Key: "", Value: "See more at http://google.com/"})

	this.assertParseSuccess("\tSee more at http://google.com/", // tab prefixed value-only line w/ colon
		&LineItem{Type: valueOnly, Key: "", Value: "See more at http://google.com/"})

	this.assertParseSuccess("", // separator lines
		&LineItem{Type: separator})

	this.assertParseSuccess("#comment line", // comment lines
		&LineItem{Type: comment, Value: "#comment line"})

	this.assertParseFailure("Invalid Key: value")
}

func (this *LineItemFixture) assertParseFailure(input string) {
	parsed, err := parse(input)
	this.So(parsed, should.BeNil)
	this.So(err, should.NotBeNil)
}
