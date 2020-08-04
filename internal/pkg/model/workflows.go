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
	"github.com/arnumina/uuid"
)

func (m *Model) stepToJob(wf *jw.Workflow, stepName string) (*jw.Job, error) {
	step, ok := wf.Steps[stepName]
	if !ok {
		return nil,
			failure.New(nil).
				Set("step", stepName).
				Msg("this step does not exist") ////////////////////////////////////////////////////////////////////////
	}

	job := jw.NewJob(
		uuid.New(),
		stepName,
		step.Plugin,
		step.Type,
		wf.Origin,
		wf.Priority,
		nil,
		wf.Emails,
	)

	if step.Config != nil {
		c := make(map[string]interface{})

		for k, v := range step.Config {
			c[k] = v
		}

		job.Data["__config"] = c
	}

	job.WfID = &wf.ID
	job.Reference = wf.CreatedAt

	return job, nil
}

func (m *Model) firstJob(wf *jw.Workflow) (*jw.Job, error) {
	job, err := m.stepToJob(wf, wf.FirstStep)
	if err != nil {
		return nil, err
	}

	for k, v := range wf.Data {
		job.Data[k] = v
	}

	wfFailed := false
	job.WfFailed = &wfFailed

	return job, nil
}

func (m *Model) newWorkflow(wf *jw.Workflow) {
	m.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"New workflow [begin]",
		"id", wf.ID,
		"name", wf.Name,
		"title", wf.Title,
		"origin", wf.Origin,
		"priority", wf.Priority,
	)

	m.channel <- message.New("new.workflow", wf)
}

// InsertWorkflow AFAIRE.
func (m *Model) InsertWorkflow(wf *jw.Workflow) error {
	job, err := m.firstJob(wf)
	if err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Impossible to determine the first job in this workflow",
			"id", wf.ID,
			"title", wf.Title,
			"reason", err.Error(),
		)

		return err
	}

	if err := m.backend.InsertWorkflow(wf, job); err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Impossible to insert a new workflow",
			"id", wf.ID,
			"title", wf.Title,
			"reason", err.Error(),
		)

		return err
	}

	m.newWorkflow(wf)
	m.newJob(job)

	return nil
}

func (m *Model) workflow(id string) (*jw.Workflow, error) {
	wf, err := m.backend.Workflow(id)
	if err != nil {
		return nil, err
	}

	if wf == nil {
		return nil,
			failure.New(nil).
				Set("workflow", id).
				Msg("this workflow does not seem to exist") ////////////////////////////////////////////////////////////
	}

	return wf, nil
}

func (m *Model) nextStep(job *jw.Job, wf *jw.Workflow) (string, error) {
	step, ok := wf.Steps[job.Name]
	if !ok {
		return "",
			failure.New(nil).
				Set("step", job.Name).
				Msg("this step does not exist") ////////////////////////////////////////////////////////////////////////
	}

	if step.Next == nil {
		return "", nil
	}

	if job.NextStep != nil {
		return *job.NextStep, nil
	}

	names := []string{}

	if job.Result != nil {
		names = append(names, *job.Result)
	}

	names = append(names, string(job.Status), "default")

	for _, r := range names {
		next, ok := step.Next[r]
		if !ok {
			continue
		}

		switch nv := next.(type) {
		case nil:
			return "", nil
		case string:
			return nv, nil
		default:
			return "", failure.Unexpected() // AFINIR
		}
	}

	return "",
		failure.New(nil).
			Set("job", job.ID).
			Set("steps", names).
			Msg("impossible to determine the next step") ///////////////////////////////////////////////////////////////
}

func (m *Model) workflowFinished(wf *jw.Workflow) {
	m.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"Workflow finished",
		"id", wf.ID,
		"name", wf.Name,
		"title", wf.Title,
		"status", wf.Status,
	)

	m.channel <- message.New("workflow.ended", wf)
}

func (m *Model) nextJob(wf *jw.Workflow, pJob *jw.Job, stepName string) (*jw.Job, error) {
	job, err := m.stepToJob(wf, stepName)
	if err != nil {
		return nil, err
	}

	for k, v := range pJob.Data {
		if k != "" && k[0] != '_' {
			job.Data[k] = v
		}
	}

	job.WfFailed = pJob.WfFailed

	return job, nil
}

func (m *Model) updateWorkflow(job *jw.Job, wf *jw.Workflow) error {
	step, err := m.nextStep(job, wf)
	if err != nil {
		return err
	}

	if step == "" {
		now := time.Now()

		if *job.WfFailed {
			wf.Status = jw.Failed
		} else {
			wf.Status = jw.Succeeded
		}

		wf.FinishedAt = &now

		if err := m.backend.UpdateWorkflow(wf); err != nil {
			return err
		}

		m.workflowFinished(wf)

		return nil
	}

	job, err = m.nextJob(wf, job, step)
	if err != nil {
		return err
	}

	return m.InsertJob(job)
}

/*
######################################################################################################## @(°_°)@ #######
*/
