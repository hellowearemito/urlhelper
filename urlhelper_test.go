package urlhelper

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURLHelper(t *testing.T) {
	Convey("URL Helper", t, func() {
		assert := func(expected, actual string) {
			So(actual, ShouldEqual, expected)
		}

		Convey("With no base url, no scheme in host and no :// in default scheme", func() {
			h := New("example.org", "", "http")

			assert("/test/url", h.Relative("test/url"))
			assert("/test/url", h.Relative("/test/url"))
			assert("/test/url/", h.Relative("/test/url/"))
			assert("/test/url?foo=bar", h.Relative("test/url", SimpleParams{"foo": "bar"}))
			assert("/test/url?foo=bar%2Fbaz", h.Relative("test/url", SimpleParams{"foo": "bar/baz"}))
			assert("http://example.org/test/url?foo=bar%2Fbaz", h.Absolute("test/url", SimpleParams{"foo": "bar/baz"}))
			assert("://example.org/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "://", SimpleParams{"foo": "bar/baz"}))
			assert("ftp://example.org/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "ftp", SimpleParams{"foo": "bar/baz"}))
			assert("ftp://example.org/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "ftp://", SimpleParams{"foo": "bar/baz"}))

			// test url.Values
			assert("/test/url?foo=bar&foo=baz", h.Relative("test/url", url.Values{"foo": {"bar", "baz"}}))
		})

		Convey("With base url with trailing slash and no leading slash, scheme in host, host with leading slash, :// in default scheme, schemes don't match", func() {
			h := New("http://example.org/", "base/url/", "https://")

			assert("/base/url/test/url", h.Relative("test/url"))
			assert("/base/url/test/url", h.Relative("/test/url"))
			assert("/base/url/test/url/", h.Relative("/test/url/"))
			assert("/base/url/test/url?foo=bar", h.Relative("test/url", SimpleParams{"foo": "bar"}))
			assert("/base/url/test/url?foo=bar%2Fbaz", h.Relative("test/url", SimpleParams{"foo": "bar/baz"}))
			assert("https://example.org/base/url/test/url?foo=bar%2Fbaz", h.Absolute("test/url", SimpleParams{"foo": "bar/baz"}))
			assert("://example.org/base/url/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "://", SimpleParams{"foo": "bar/baz"}))
			assert("ftp://example.org/base/url/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "ftp", SimpleParams{"foo": "bar/baz"}))
			assert("ftp://example.org/base/url/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "ftp://", SimpleParams{"foo": "bar/baz"}))
		})

		Convey("With protocol-relative default scheme, host with different scheme, and base url with leading slash, but not trailing one", func() {
			h := New("http://example.org/", "/base/url", "://")

			assert("/base/url/test/url", h.Relative("test/url"))
			assert("/base/url/test/url", h.Relative("/test/url"))
			assert("/base/url/test/url/", h.Relative("/test/url/"))
			assert("/base/url/test/url?foo=bar", h.Relative("test/url", SimpleParams{"foo": "bar"}))
			assert("/base/url/test/url?foo=bar%2Fbaz", h.Relative("test/url", SimpleParams{"foo": "bar/baz"}))
			assert("://example.org/base/url/test/url?foo=bar%2Fbaz", h.Absolute("test/url", SimpleParams{"foo": "bar/baz"}))
			assert("://example.org/base/url/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "://", SimpleParams{"foo": "bar/baz"}))
			assert("ftp://example.org/base/url/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "ftp", SimpleParams{"foo": "bar/baz"}))
			assert("ftp://example.org/base/url/test/url?foo=bar%2Fbaz", h.Scheme("test/url", "ftp://", SimpleParams{"foo": "bar/baz"}))
		})

		Convey("Panics", func() {
			h := New("http://example.org/", "base/url/", "https://")

			params := []Params{url.Values{
				"foo": {"bar"},
			}, url.Values{
				"baz": {"qux"},
			}}

			So(func() {
				h.Relative("", params...)
			}, ShouldPanic)
			So(func() {
				h.Absolute("", params...)
			}, ShouldPanic)
			So(func() {
				h.Scheme("", "://", params...)
			}, ShouldPanic)
		})
	})
}
