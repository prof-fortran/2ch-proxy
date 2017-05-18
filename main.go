package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func getPortAndHostFromEnvironment() (string, string) {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	host := os.Getenv("2CH_HOST")

	if host == "" {
		log.Fatal("$HOST must be set")
	}
	return port, host
}

func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

func main() {
	port, host := getPortAndHostFromEnvironment()
	log.Printf(
		"Starting application on port: %s, reversing host %s\n",
		port,
		host,
	)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   host,
	})
	singleHosted := sameHost(proxy)
	http.ListenAndServe(":"+port, singleHosted)
}
