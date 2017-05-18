package main

import (
	"log"
	"net/url"
	"os"
)

func getPortAndUrlFromEnvironment() (string, *url.URL) {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	rawUrl := os.Getenv("2CH_URL")

	if rawUrl == "" {
		log.Fatal("$2CH_URL must be set")
	}
	finalUrl, err := url.Parse(rawUrl)
	if err != nil {
			log.Fatal(err)
  }

	return port, finalUrl
}

func main() {
	port, proxyUrl := getPortAndUrlFromEnvironment()
	log.Printf(
		"Starting application on port: %s, reversing url %s\n",
		port,
		proxyUrl,
	)
	proxyServer := &ProxyServer{
		Port: port,
		ProxificationHost: proxyUrl.Hostname(),
		ProxificationScheme: proxyUrl.Scheme,
	}
	proxyServer.Start()
}
