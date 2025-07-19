package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	overrideMIME = "application/x-apple-aspen-config"
	listenAddr   = ":8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b64 := r.URL.Path[1:] // strip leading "/"

		if b64 == "" {
			http.Error(w, "Missing base64 target URL", http.StatusBadRequest)
			return
		}

		rawBytes, err := base64.URLEncoding.DecodeString(b64)
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusBadRequest)
			return
		}

		target := string(rawBytes)
		targetURL, err := url.Parse(target)
		if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
			http.Error(w, "Malformed target URL", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(target)
		if err != nil {
			http.Error(w, "Failed to fetch target: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", overrideMIME)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("Proxy running on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
