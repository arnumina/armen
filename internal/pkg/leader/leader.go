/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package leader

import (
	"sync"
	"time"

	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"

	"github.com/arnumina/armen/internal/pkg/application"
	"github.com/arnumina/armen/internal/pkg/backend"
)

const (
	_expiration = 30 * time.Second
	_retryToBe  = 20 * time.Second
)

type (
	// Resource AFAIRE.
	Resource interface {
		resources.Leader
	}

	// Leader AFAIRE.
	Leader struct {
		app     application.Resource
		logger  *logger.Logger
		backend backend.Resource
		success bool
		mutex   sync.RWMutex
		stop    chan struct{}
		stopped chan struct{}
	}
)

// New AFAIRE.
func New(app application.Resource, logger *logger.Logger, backend backend.Resource) *Leader {
	return &Leader{
		app:     app,
		logger:  logger,
		backend: backend,
		stop:    make(chan struct{}, 1),
		stopped: make(chan struct{}, 1),
	}
}

func (l *Leader) tryToBe() error {
	v, err := l.backend.Lock("leader", l.app.ID(), _expiration)

	l.mutex.Lock()

	if l.success != v {
		l.success = v
		l.logger.Notice("leader?", "success", v) //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
	}

	l.mutex.Unlock()

	return err
}

// Build AFAIRE.
func (l *Leader) Build() (*Leader, error) {
	if err := l.tryToBe(); err != nil {
		return nil,
			failure.New(err).Msg("leader") /////////////////////////////////////////////////////////////////////////////
	}

	go func() {
		for {
			select {
			case <-l.stop:
				close(l.stopped)
				return
			case <-time.After(_retryToBe):
				_ = l.tryToBe()
			}
		}
	}()

	return l, nil
}

// Success AFAIRE.
func (l *Leader) Success() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.success
}

// Close AFAIRE.
func (l *Leader) Close() {
	close(l.stop)
	<-l.stopped

	if !l.success {
		return
	}

	if err := l.backend.Unlock("leader", l.app.ID()); err != nil {
		l.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Leader.Close()",
			"reason", err.Error(),
		)
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
