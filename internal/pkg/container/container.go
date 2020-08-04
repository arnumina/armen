/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package container

import (
	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/logger"
)

type (
	// Container AFAIRE.
	Container struct {
		util   resources.Util
		app    resources.Application
		logger *logger.Logger
		bus    resources.Bus
		leader resources.Leader
		model  resources.Model
		server resources.Server
	}
)

// New AFAIRE.
func New() *Container {
	return &Container{}
}

// Util AFAIRE.
func (c *Container) Util() resources.Util {
	return c.util
}

// SetUtil AFAIRE.
func (c *Container) SetUtil(util resources.Util) {
	c.util = util
}

// Application AFAIRE.
func (c *Container) Application() resources.Application {
	return c.app
}

// SetApplication AFAIRE.
func (c *Container) SetApplication(app resources.Application) {
	c.app = app
}

// Logger AFAIRE.
func (c *Container) Logger() *logger.Logger {
	return c.logger
}

// SetLogger AFAIRE.
func (c *Container) SetLogger(logger *logger.Logger) {
	c.logger = logger
}

// Bus AFAIRE.
func (c *Container) Bus() resources.Bus {
	return c.bus
}

// SetBus AFAIRE.
func (c *Container) SetBus(bus resources.Bus) {
	c.bus = bus
}

// Leader AFAIRE.
func (c *Container) Leader() resources.Leader {
	return c.leader
}

// SetLeader AFAIRE.
func (c *Container) SetLeader(leader resources.Leader) {
	c.leader = leader
}

// Model AFAIRE.
func (c *Container) Model() resources.Model {
	return c.model
}

// SetModel AFAIRE.
func (c *Container) SetModel(model resources.Model) {
	c.model = model
}

// Server AFAIRE.
func (c *Container) Server() resources.Server {
	return c.server
}

// SetServer AFAIRE.
func (c *Container) SetServer(server resources.Server) {
	c.server = server
}

/*
######################################################################################################## @(°_°)@ #######
*/
