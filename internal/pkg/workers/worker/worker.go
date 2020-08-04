/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package worker

import (
	"time"

	"github.com/arnumina/armen.core/pkg/jw"
	"github.com/arnumina/armen.core/pkg/message"
	"github.com/arnumina/logger"

	"github.com/arnumina/armen/internal/pkg/bus"
	"github.com/arnumina/armen/internal/pkg/model"
	"github.com/arnumina/armen/internal/pkg/plugins"
	"github.com/arnumina/armen/internal/pkg/util"
	"github.com/arnumina/armen/internal/pkg/workers/runner"
)

type (
	// Worker AFAIRE.
	Worker struct {
		util       util.Resource
		model      model.Resource
		plugins    plugins.Resource
		stop       <-chan struct{}
		temporary  bool
		logger     *logger.Logger
		jobCounter int
		channel    chan<- *message.Message
	}
)

// New AFAIRE.
func New(util util.Resource, bus bus.Resource, model model.Resource, plugins plugins.Resource, stop <-chan struct{},
	temporary bool, logger *logger.Logger) *Worker {
	return &Worker{
		util:      util,
		model:     model,
		plugins:   plugins,
		stop:      stop,
		temporary: temporary,
		logger:    logger,
		channel:   bus.AddPublisher("worker", 1, 1),
	}
}

func (w *Worker) newLogger(job *jw.Job) *logger.Logger {
	if job.WfID == nil {
		return w.logger.Clone(w.util.LoggerPrefix("job", job.ID[:8]))
	}

	return w.logger.Clone(w.util.LoggerPrefix("workflow", (*job.WfID)[:8]))
}

func (w *Worker) runJob(job *jw.Job) {
	w.jobCounter++

	w.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"Run job",
		"id", job.ID,
		"count", w.jobCounter,
	)

	runner.New(w.model, w.newLogger(job), job, w.channel).Run(w.plugins.Find(job.Plugin))
}

// Run AFAIRE.
func (w *Worker) Run() {
	var delay time.Duration

loop:
	for {
		select {
		case <-w.stop:
			break loop
		case <-time.After(delay * time.Second):
			if job := w.model.NextJob(); job != nil {
				w.runJob(job)

				delay = 0
			} else {
				if w.temporary {
					break loop
				}

				delay = 1
			}
		}
	}

	close(w.channel)
}

/*
######################################################################################################## @(°_°)@ #######
*/
