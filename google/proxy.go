package google

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

// The proxy return a http.Handler that will proxy all requests to the Google Fonts server
// it take a path, which is used when generating the urls of the font bytes
// For example, if you want to proxy requests on `example.com/fonts`:
//
//     http.Handle("/fonts", http.StripPrefix("/fonts", google.Proxy("/fonts")))
func Proxy(path string) http.Handler {
	apiURL, _ := url.Parse("https://fonts.googleapis.com")
	staticURL, _ := url.Parse("https://fonts.gstatic.com")

	director := func(req *http.Request) {
		var target *url.URL
		switch {
		case strings.HasPrefix(req.URL.Path, "/static"):
			target = staticURL
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/static")
			req.URL.RawPath = strings.TrimPrefix(req.URL.RawPath, "/static")
		default:
			// so we can replace the gstatic url without having to gzip decode
			req.Header.Del("Accept-Encoding")
			target = apiURL
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{Director: director, ModifyResponse: responseModifier(path)}
}

func responseModifier(path string) func(resp *http.Response) error {
	return func(resp *http.Response) error {
		oldBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		replaced := bytes.ReplaceAll(oldBody,
			[]byte("https://fonts.gstatic.com"),
			[]byte(filepath.ToSlash(filepath.Join(path, "static"))),
		)

		resp.Body = io.NopCloser(bytes.NewBuffer(replaced))
		return nil
	}
}
