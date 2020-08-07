/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package config

import (
	"github.com/arnumina/config"
	"github.com/arnumina/failure"

	"github.com/arnumina/armen/internal/pkg/application"
	"github.com/arnumina/armen/internal/pkg/util"
)

type (
	// Resource AFAIRE.
	Resource interface {
		Backend() *Backend
		Logger() *Logger
		Server() *Server
		Workers() *Workers
	}

	// Backend AFAIRE.
	Backend struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}

	// Logger AFAIRE.
	Logger struct {
		Level     string
		Formatter string
		Output    string
		Syslog    struct {
			Facility string
		}
	}

	// Server AFAIRE.
	Server struct {
		Port     int
		TLS      bool
		CertFile string
		KeyFile  string
	}

	// Workers AFAIRE.
	Workers struct {
		Count int
	}

	// Config AFAIRE.
	Config struct {
		util      util.Resource
		app       application.Resource
		version   bool
		cfgString string
		key       string
		decrypt   string
		encrypt   string
		port      int
		cfg       struct {
			Backend *Backend
			Logger  *Logger
			Server  *Server
			Workers *Workers
		}
	}
)

// New AFAIRE.
func New(util util.Resource, app application.Resource) *Config {
	return &Config{
		util:      util,
		app:       app,
		cfgString: "empty",
	}
}

func (c *Config) load() error {
	if err := c.defaultConfigString(); err != nil {
		return err
	}

	if err := c.parseFlag(); err != nil {
		return err
	}

	if err := c.handleFlag(); err != nil {
		return err
	}

	cfg, err := config.Load(c.cfgString)
	if err != nil {
		return err
	}

	if err := c.util.DecodeData(cfg, &c.cfg); err != nil {
		return err
	}

	if c.port != 0 {
		c.cfg.Server.Port = c.port
	}

	return nil
}

// Load AFAIRE.
func (c *Config) Load() (*Config, error) {
	if err := c.load(); err != nil {
		return nil,
			failure.New(err).Msg("config") /////////////////////////////////////////////////////////////////////////////
	}

	return c, nil
}

// Backend AFAIRE.
func (c *Config) Backend() *Backend {
	return c.cfg.Backend
}

// Logger AFAIRE.
func (c *Config) Logger() *Logger {
	return c.cfg.Logger
}

// Server AFAIRE.
func (c *Config) Server() *Server {
	return c.cfg.Server
}

// Workers AFAIRE.
func (c *Config) Workers() *Workers {
	return c.cfg.Workers
}

/*
######################################################################################################## @(°_°)@ #######
*/
