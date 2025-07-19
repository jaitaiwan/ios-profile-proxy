package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	targetBase     = "https://ios.a.proxy.q3k.onl"      // Replace with your backend
	overrideMIME   = "application/x-apple-aspen-config" // Replace with desired MIME type
	listenAddr     = ":8080"
)

func main() {
	targetURL, err := url.Parse(targetBase)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			originalPath := req.URL.Path
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = singleJoiningSlash(targetURL.Path, originalPath)
			if targetURL.RawQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetURL.RawQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetURL.RawQuery + "&" + req.URL.RawQuery
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Set("Content-Type", overrideMIME)
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "Proxy error: "+err.Error(), http.StatusBadGateway)
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Proxy listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	default:
		return a + b
	}
}
