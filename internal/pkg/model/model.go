/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package model

import (
	"time"

	"github.com/arnumina/armen.core/pkg/jw"
	"github.com/arnumina/armen.core/pkg/message"
	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"

	"github.com/arnumina/armen/internal/pkg/backend"
	"github.com/arnumina/armen/internal/pkg/bus"
)

type (
	// Resource AFAIRE.
	Resource interface {
		resources.Model
		NextJob() *jw.Job
		UpdateJob(job *jw.Job)
	}

	// Model AFAIRE.
	Model struct {
		logger  *logger.Logger
		bus     bus.Resource
		backend backend.Resource
		channel chan<- *message.Message
		njLimit time.Time
	}
)

// New AFAIRE.
func New(logger *logger.Logger, bus bus.Resource, backend backend.Resource) *Model {
	return &Model{
		logger:  logger,
		bus:     bus,
		backend: backend,
		channel: bus.AddPublisher("model", 1, 1),
		njLimit: time.Now(),
	}
}

// Build AFAIRE.
func (m *Model) Build() (*Model, error) {
	if err := m.subscribe(); err != nil {
		m.Close()
		return nil,
			failure.New(err).Msg("model") //////////////////////////////////////////////////////////////////////////////
	}

	return m, nil
}

// Close AFAIRE.
func (m *Model) Close() {
	close(m.channel)
}

/*
######################################################################################################## @(°_°)@ #######
*/
