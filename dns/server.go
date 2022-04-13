package dns

import (
	"fmt"

	"gozar/config"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

type Server struct {
	srv    *dns.Server
	Config *config.Config
}

func New(
	c *config.Config,
) *Server {
	srv := &dns.Server{Addr: fmt.Sprintf("%s:%d", c.DNS.Host, c.DNS.Port), Net: "udp"}

	return &Server{
		Config: c,
		srv:    srv,
	}
}

func (s *Server) Start() error {
	logrus.Infof("Starting dns server at: %s:%d", s.Config.DNS.Host, s.Config.DNS.Port)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.srv.Shutdown()
}

func (s *Server) SetHandler(handler dns.Handler) {
	s.srv.Handler = handler
}
