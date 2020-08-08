/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"
	"github.com/gorilla/mux"

	"github.com/arnumina/armen/internal/pkg/config"
	"github.com/arnumina/armen/internal/pkg/util"
)

type (
	// Resource AFAIRE.
	Resource interface {
		resources.Server
	}

	// Server AFAIRE.
	Server struct {
		logger   *logger.Logger
		adapter  *logger.LogAdapter
		port     int
		tls      bool
		certFile string
		keyFile  string
		server   *http.Server
		router   *mux.Router
		stopped  chan error
	}
)

// New AFAIRE.
func New(util util.Resource, logger *logger.Logger, config config.Resource) *Server {
	logger = util.CloneLogger(logger, "server")
	cfg := config.Server()
	adapter := logger.NewLogAdapter("error")
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     adapter.NewLogger("", log.Llongfile),
	}

	if cfg.TLS {
		server.TLSConfig = &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		}
	}

	return &Server{
		logger:   logger,
		adapter:  adapter,
		port:     cfg.Port,
		tls:      cfg.TLS,
		certFile: cfg.CertFile,
		keyFile:  cfg.KeyFile,
		server:   server,
		router:   router,
		stopped:  make(chan error, 1),
	}
}

func (s *Server) start() error {
	go func() {
		var err error

		if s.tls {
			err = s.server.ListenAndServeTLS(s.certFile, s.keyFile)
		} else {
			err = s.server.ListenAndServe()
		}

		s.stopped <- err
		close(s.stopped)
	}()

	select {
	case err := <-s.stopped:
		return err
	case <-time.After(250 * time.Millisecond):
		s.logger.Info(">>>Server", "port", s.port) //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		return nil
	}
}

// Start AFAIRE.
func (s *Server) Start() (*Server, error) {
	if err := s.start(); err != nil {
		return nil,
			failure.New(err).Msg("server") /////////////////////////////////////////////////////////////////////////////
	}

	return s, nil
}

// Port AFAIRE.
func (s *Server) Port() int {
	return s.port
}

// Router AFAIRE.
func (s *Server) Router() *mux.Router {
	return s.router
}

// Stop AFAIRE.
func (s *Server) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Server.Stop() - Shutdown()",
			"reason", err.Error(),
		)
	}

	if err := <-s.stopped; !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Server.Stop() - ListenAndServe[TLS]()",
			"reason", err.Error(),
		)
	}

	<-s.stopped

	s.logger.Info("<<<Server") //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
}

/*
######################################################################################################## @(°_°)@ #######
*/
