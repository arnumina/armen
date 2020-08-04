/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package backend

import (
	"errors"

	"github.com/arnumina/armen.core/pkg/jw"
	"github.com/jackc/pgx/v4"
)

// InsertWorkflow AFAIRE.
func (b *Backend) InsertWorkflow(wf *jw.Workflow, job *jw.Job) error {
	t, err := b.pgc.Begin()
	if err != nil {
		return err
	}

	defer t.Rollback()

	_, err = t.Exec(
		`INSERT INTO workflows (id, name, title, origin, priority,
		first_step, steps, emails, data, status, created_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		wf.ID,
		wf.Name,
		wf.Title,
		wf.Origin,
		wf.Priority,
		wf.FirstStep,
		wf.Steps,
		wf.Emails,
		wf.Data,
		wf.Status,
		wf.CreatedAt,
		wf.FinishedAt,
	)
	if err != nil {
		return err
	}

	_, err = t.Exec(jobSQLInsert(), jobColumns(job)...)
	if err != nil {
		return err
	}

	return t.Commit()
}

// Workflow AFAIRE.
func (b *Backend) Workflow(id string) (*jw.Workflow, error) {
	var wf jw.Workflow

	if err := b.pgc.QueryRow(
		`SELECT id, name, title, origin, priority, first_step, steps, emails, data, status, created_at, finished_at
		FROM workflows
		WHERE id = $1`,
		id,
	).Scan(
		&wf.ID,
		&wf.Name,
		&wf.Title,
		&wf.Origin,
		&wf.Priority,
		&wf.FirstStep,
		&wf.Steps,
		&wf.Emails,
		&wf.Data,
		&wf.Status,
		&wf.CreatedAt,
		&wf.FinishedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &wf, nil
}

// UpdateWorkflow AFAIRE.
func (b *Backend) UpdateWorkflow(wf *jw.Workflow) error {
	_, err := b.pgc.Exec(
		"UPDATE workflows SET status = $1, finished_at = $2 WHERE id = $3",
		wf.Status,
		wf.FinishedAt,
		wf.ID,
	)

	return err
}

/*
######################################################################################################## @(°_°)@ #######
*/
