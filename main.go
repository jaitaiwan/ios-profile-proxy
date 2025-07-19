package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	overrideMIME = "application/x-apple-aspen-config"
	listenAddr   = ":8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b64 := strings.TrimPrefix(r.URL.Path, "/")
		if b64 == "" {
			http.Error(w, "Missing base64 target URL", http.StatusBadRequest)
			return
		}

		rawURLBytes, err := base64.URLEncoding.DecodeString(b64)
		if err != nil {
			http.Error(w, "Invalid base64 URL", http.StatusBadRequest)
			return
		}
		targetURL := string(rawURLBytes)

		parsedURL, err := url.Parse(targetURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			http.Error(w, "Malformed target URL", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, "Failed to fetch target: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Extract file name from URL path
		filename := path.Base(parsedURL.Path)
		if filename == "" || filename == "/" {
			filename = "download.mobileconfig"
		}

		// Set headers
		w.Header().Set("Content-Type", overrideMIME)
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
	})

	log.Printf("Proxy listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
