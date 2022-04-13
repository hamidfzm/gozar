package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"gozar/config"
	"gozar/proxy/handlers"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
)

type Server struct {
	srvs   []*http.Server
	Config *config.Config
}

func New(
	c *config.Config,
) *Server {
	var servers []*http.Server
	for _, port := range c.Proxy.Ports {
		srv := &http.Server{
			Addr: fmt.Sprintf("%s:%d", c.Proxy.Host, port),
			// Disable HTTP/2.
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}

		if port == 443 {
			// srv.TLSConfig = m.TLSConfig()
		}

		srv.Handler = handlers.New(srv)
		servers = append(servers, srv)
	}

	return &Server{
		Config: c,
		srvs:   servers,
	}
}

func (s *Server) GoListenAndServe(index int) {
	err := s.srvs[index].ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (s *Server) Start() error {
	for i := range s.srvs {
		logrus.Infof("Starting proxy server at: %s", s.srvs[i].Addr)
		go s.GoListenAndServe(i)
	}
	return nil
}

func (s *Server) Shutdown() error {
	background := context.Background()
	var result error
	for i := range s.srvs {
		err := s.srvs[i].Shutdown(background)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}
