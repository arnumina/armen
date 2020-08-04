/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package runner

import (
	"time"

	"github.com/arnumina/armen.core/pkg/jw"
	"github.com/arnumina/armen.core/pkg/message"
	"github.com/arnumina/armen.core/pkg/plugin"
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"

	"github.com/arnumina/armen/internal/pkg/model"
)

type (
	// Runner AFAIRE.
	Runner struct {
		model   model.Resource
		logger  *logger.Logger
		job     *jw.Job
		channel chan<- *message.Message
	}
)

// New AFAIRE.
func New(model model.Resource, logger *logger.Logger, job *jw.Job, ch chan<- *message.Message) *Runner {
	return &Runner{
		model:   model,
		logger:  logger,
		job:     job,
		channel: ch,
	}
}

func (r *Runner) succeeded() {
	r.job.Status = jw.Succeeded
}

func (r *Runner) pending(jwr *jw.Result) {
	r.job.Status = jw.Pending
	r.job.RunAfter = time.Now().Add(jwr.Delay)

	if jwr.Failure != nil {
		r.job.Attempt++

		if r.job.Attempt >= 5 {
			r.logger.Notice( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
				"The number of attempts to execute this job is high",
				"id", r.job.ID,
				"plugin", r.job.Plugin,
				"type", r.job.Type,
				"priority", r.job.Priority,
				"Attempt", r.job.Attempt,
			)
		}
	}
}

func (r *Runner) failed() {
	r.job.Status = jw.Failed
}

// Run AFAIRE.
func (r *Runner) Run(plugin plugin.Plugin) {
	if r.job.Status == jw.Todo {
		r.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Begin",
			"name", r.job.Name,
			"plugin", r.job.Plugin,
			"type", r.job.Type,
		)
	} else {
		r.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Resume",
			"name", r.job.Name,
			"plugin", r.job.Plugin,
			"type", r.job.Type,
			"attempt", r.job.Attempt,
		)
	}

	r.job.Status = jw.Running

	r.channel <- message.New("job.before.run", r.job)

	var jwr *jw.Result

	if plugin == nil {
		jwr = jw.NewResult(
			failure.New(nil).Msg("the application specified for this job does not exist"), /////////////////////////////
			0,
		)
	} else {
		jwr = plugin.RunJob(r.job, r.logger)
	}

	if jwr == nil {
		r.succeeded()
	} else if jwr.Failure == nil {
		r.pending(jwr)
	} else {
		r.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"The execution of this job didn't work",
			"id", r.job.ID,
			"plugin", r.job.Plugin,
			"type", r.job.Type,
			"priority", r.job.Priority,
			"attempt", r.job.Attempt,
			"reason", jwr.Failure.Error(),
		)

		if jwr.Delay == 0 {
			r.failed()
		} else {
			r.pending(jwr)
		}

		r.job.Data["__error"] = jwr.Failure.Error()
	}

	r.channel <- message.New("job.after.run", r.job)

	if r.job.Status == jw.Pending {
		r.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Continuation",
			"after", r.job.RunAfter.Round(time.Second).String(),
			"attempt", r.job.Attempt,
		)
	} else {
		r.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"End",
			"status", r.job.Status,
		)
	}

	r.model.UpdateJob(r.job)
}

/*
######################################################################################################## @(°_°)@ #######
*/
