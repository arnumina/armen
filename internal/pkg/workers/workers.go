/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package workers

import (
	"sync"

	"github.com/arnumina/logger"
	"github.com/arnumina/uuid"

	"github.com/arnumina/armen/internal/pkg/bus"
	"github.com/arnumina/armen/internal/pkg/config"
	"github.com/arnumina/armen/internal/pkg/model"
	"github.com/arnumina/armen/internal/pkg/plugins"
	"github.com/arnumina/armen/internal/pkg/util"
	"github.com/arnumina/armen/internal/pkg/workers/worker"
)

type (
	// Workers AFAIRE.
	Workers struct {
		util    util.Resource
		logger  *logger.Logger
		bus     bus.Resource
		model   model.Resource
		plugins plugins.Resource
		group   sync.WaitGroup
		stop    chan struct{}
	}
)

// New AFAIRE.
func New(util util.Resource, logger *logger.Logger, bus bus.Resource, model model.Resource,
	plugins plugins.Resource) *Workers {
	return &Workers{
		util:    util,
		logger:  logger,
		bus:     bus,
		model:   model,
		plugins: plugins,
		stop:    make(chan struct{}),
	}
}

func (w *Workers) goWorker(temporary bool) {
	go func() {
		logger := w.logger.Clone(w.util.LoggerPrefix("worker", uuid.New()))

		if temporary {
			logger.Info("+++Worker") //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		} else {
			logger.Info(">>>Worker") //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		}

		worker.New(w.util, w.bus, w.model, w.plugins, w.stop, temporary, logger).Run()

		if temporary {
			logger.Info("---Worker") //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		} else {
			logger.Info("<<<Worker") //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		}

		w.group.Done()
	}()
}

func (w *Workers) add(temporary bool) {
	w.group.Add(1)
	w.goWorker(temporary)
}

// Start AFAIRE.
func (w *Workers) Start(config config.Resource) *Workers {
	cfg := config.Workers()

	for n := 0; n < cfg.Count; n++ {
		w.add(false)
	}

	return w
}

// Add AFAIRE.
func (w *Workers) Add() {
	w.add(true)
}

// Stop AFAIRE.
func (w *Workers) Stop() {
	close(w.stop)
	w.group.Wait()
}

/*
######################################################################################################## @(°_°)@ #######
*/
