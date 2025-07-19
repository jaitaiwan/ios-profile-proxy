package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	overrideMIME = "application/x-apple-aspen-config" // typical for .mobileconfig
	listenAddr   = ":8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetRaw := strings.TrimPrefix(r.URL.Path, "/")
		if targetRaw == "" {
			http.Error(w, "Missing target URL in path", http.StatusBadRequest)
			return
		}

		targetURL, err := url.QueryUnescape(targetRaw)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.Parse(targetURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			http.Error(w, "Malformed target URL", http.StatusBadRequest)
			return
		}

		// Fetch the content
		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, "Error fetching target: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Overwrite MIME type
		w.Header().Set("Content-Type", overrideMIME)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("Proxy server listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
