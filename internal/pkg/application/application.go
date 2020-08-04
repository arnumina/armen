/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package application

import (
	"time"

	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/failure"
	"github.com/arnumina/uuid"
)

type (
	// Resource AFAIRE.
	Resource interface {
		resources.Application
		Build() error
	}

	// Application AFAIRE.
	Application struct {
		id        string
		name      string
		version   string
		builtAt   time.Time
		startedAt time.Time
		fqdn      string
	}
)

// New AFAIRE.
func New(name, version string, builtAt time.Time) *Application {
	return &Application{
		id:        uuid.New(),
		name:      name,
		version:   version,
		builtAt:   builtAt,
		startedAt: time.Now(),
	}
}

// ID AFAIRE.
func (a *Application) ID() string {
	return a.id
}

// Name AFAIRE.
func (a *Application) Name() string {
	return a.name
}

// Version AFAIRE.
func (a *Application) Version() string {
	return a.version
}

// BuiltAt AFAIRE.
func (a *Application) BuiltAt() time.Time {
	return a.builtAt
}

// StartedAt AFAIRE.
func (a *Application) StartedAt() time.Time {
	return a.startedAt
}

// FQDN AFAIRE.
func (a *Application) FQDN() string {
	return a.fqdn
}

// Build AFAIRE.
func (a *Application) Build() error {
	fqdn, err := getFQDN()
	if err != nil {
		return failure.New(err).Msg("application") /////////////////////////////////////////////////////////////////////
	}

	a.fqdn = fqdn

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
