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
	"github.com/arnumina/failure"
)

func (m *Model) newJob(job *jw.Job) {
	var wf string

	if job.WfID == nil {
		wf = "nil"
	} else {
		wf = *job.WfID
	}

	m.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"New job",
		"id", job.ID,
		"name", job.Name,
		"plugin", job.Plugin,
		"type", job.Type,
		"origin", job.Origin,
		"priority", job.Priority,
		"workflow", wf,
	)

	m.channel <- message.New("new.job", job)
}

// InsertJob AFAIRE.
func (m *Model) InsertJob(job *jw.Job) error {
	if err := m.backend.InsertJob(job); err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Impossible to insert a new job",
			"id", job.ID,
			"plugin", job.Plugin,
			"type", job.Type,
			"reason", err.Error(),
		)

		return err
	}

	m.newJob(job)

	return nil
}

// MaybeInsertJob AFAIRE.
func (m *Model) MaybeInsertJob(job *jw.Job) (bool, error) {
	if job.Key == nil {
		failure := failure.New(nil).
			Set("id", job.ID).
			Set("plugin", job.Plugin).
			Set("type", job.Type).
			Msg("the key to this job is missing") //////////////////////////////////////////////////////////////////////

		m.logger.Error(failure.Error()) //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

		return false, failure
	}

	inserted, err := m.backend.MaybeInsertJob(job)
	if err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Impossible to insert a new job",
			"id", job.ID,
			"plugin", job.Plugin,
			"type", job.Type,
			"reason", err.Error(),
		)

		return false, err
	}

	if !inserted {
		m.logger.Notice( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"A job with the same key already exists",
			"plugin", job.Plugin,
			"type", job.Type,
			"key", *job.Key,
		)

		return false, nil
	}

	m.newJob(job)

	return true, nil
}

// NextJob AFAIRE.
func (m *Model) NextJob() *jw.Job {
	if m.njLimit.After(time.Now()) {
		return nil
	}

	job, err := m.backend.NextJob()
	if err != nil {
		m.logger.Warning( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Cannot retrieve the next job to run",
			"reason", err.Error(),
		)

		m.njLimit = time.Now().Add(10 * time.Minute)

		return nil
	}

	return job
}

// UpdateJob AFAIRE.
func (m *Model) UpdateJob(job *jw.Job) {
	if err := m.backend.UpdateJob(job); err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Impossible to update this job",
			"id", job.ID,
			"plugin", job.Plugin,
			"type", job.Type,
			"reason", err.Error(),
		)

		return
	}

	if job.WfID == nil || job.Status == jw.Pending {
		return
	}

	wf, err := m.workflow(*job.WfID)
	if err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Cannot retrieve the workflow associated with this job",
			"job", *job.WfID,
			"reason", err.Error(),
		)

		return
	}

	if err := m.updateWorkflow(job, wf); err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Impossible to update this workflow",
			"id", wf.ID,
			"title", wf.Title,
			"reason", err.Error(),
		)
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
