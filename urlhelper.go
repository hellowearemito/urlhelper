package urlhelper

import (
	"net/url"
	"strings"
)

// Params is the interface for query parameters
type Params interface {
	Encode() string
}

// SimpleParams is a simplified version of url.Values that only allows
// one value per key.
type SimpleParams map[string]string

// Encode encodes the parameters into "URL encoded" form ("bar=baz&foo=quux") sorted by key.
func (p SimpleParams) Encode() string {
	values := url.Values{}

	for key, value := range p {
		values.Add(key, value)
	}

	return values.Encode()
}

// Helper is the url helper
type Helper interface {
	// Relative generates a root-relative url.
	Relative(path string, params ...Params) string
	// Absolute generates an absolute url with the default scheme
	Absolute(path string, params ...Params) string
	// Scheme generates an absolute url with the specified scheme
	Scheme(path string, scheme string, params ...Params) string
}

type helper struct {
	host    string // never has trailing slash
	baseURL string // always has leading slash, except if empty, never has trailing slash
	scheme  string
}

// New creates a new url helper
func New(host string, baseURL string, scheme string) Helper {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}

	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	if !strings.HasPrefix(baseURL, "/") && baseURL != "" {
		baseURL = "/" + baseURL
	}

	return &helper{
		host:    host,
		baseURL: baseURL,
		scheme:  scheme,
	}
}

// generate generates an url.
// path is the path relative to baseURL.
// params are the query parameters.
// if absolute is true, an absolute url will be generated, otherwise a root-relative one.
// scheme specifies the scheme for an absolute url, and is ignored if absolute is false.
// if scheme is an empty string, it will default to either http or https, depending on h.https
func (h *helper) generate(path string, params []Params, absolute bool, scheme string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := h.baseURL + path
	if len(params) > 0 {
		paramsEncoded := params[0].Encode()
		if paramsEncoded != "" {
			url = url + "?" + paramsEncoded
		}
	}

	if absolute {
		if scheme == "" {
			scheme = h.scheme
		}
		if !strings.HasSuffix(scheme, "://") {
			scheme = scheme + "://"
		}

		host := h.host
		if index := strings.Index(host, "://"); index != -1 {
			host = scheme + host[index+3:]
		} else {
			host = scheme + host
		}

		url = host + url
	}

	return url
}

// Relative generates a root-relative url.
func (h *helper) Relative(path string, params ...Params) string {
	if len(params) > 1 {
		panic("Helper.Relative called with multiple params")
	}

	return h.generate(path, params, false, "")
}

// Absolute generates an absolute url with the default scheme
func (h *helper) Absolute(path string, params ...Params) string {
	if len(params) > 1 {
		panic("Helper.Absolute called with multiple params")
	}

	return h.generate(path, params, true, "")
}

// Scheme generates an absolute url with the specified scheme
func (h *helper) Scheme(path string, scheme string, params ...Params) string {
	if len(params) > 1 {
		panic("Helper.Scheme called with multiple params")
	}

	return h.generate(path, params, true, scheme)
}
