package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyServer struct {
	Port string
	ProxificationHost string
	ProxificationScheme string
}

func (p *ProxyServer) Start(){
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: p.ProxificationScheme,
		Host:   p.ProxificationHost,
	})
	proxy.ModifyResponse = p.modifyPermanentRedirects
	configuredProxy := p.setRequestHostFromInitialRequest(
		logProxifiedRequests(proxy),
	)
	http.ListenAndServe(":"+p.Port, configuredProxy)
}

func (p *ProxyServer) modifyPermanentRedirects(response *http.Response) error {
  location, err := response.Location()
	if err == nil && location.Hostname() == p.ProxificationHost {
		response.Header.Set("Location", location.RequestURI())
	}
	return nil
}

func (p *ProxyServer) setRequestHostFromInitialRequest(
	handler http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		if r.Host == "" {
			r.Host = p.ProxificationHost
		}

		handler.ServeHTTP(w, r)
	})
}

func logProxifiedRequests(
	handler http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  log.Printf(
	  	"Proxifying request to %s\n",
			r.URL,
	  )
		handler.ServeHTTP(w, r)
	})
}
