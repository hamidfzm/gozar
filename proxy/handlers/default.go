package handlers

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type DefaultHandler struct {
	srv *http.Server
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func (h DefaultHandler) handleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusServiceUnavailable,
		)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func (h DefaultHandler) handleHTTP(w http.ResponseWriter, ireq *http.Request) {
	req := new(http.Request)
	*req = *ireq // shallow clone
	req.URL.Scheme = "http"
	req.URL.Host = ireq.Host
	req.URL.Path = ireq.URL.Path

	if h.srv.TLSConfig != nil {
		req.URL.Scheme = "https"
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (h DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Proxy[%s] %s(%s)", h.srv.Addr, r.Host, r.Method)
	if r.Method == http.MethodConnect {
		h.handleTunneling(w, r)
		return
	}
	h.handleHTTP(w, r)
}

func New(server *http.Server) http.Handler {
	return &DefaultHandler{
		srv: server,
	}
}
