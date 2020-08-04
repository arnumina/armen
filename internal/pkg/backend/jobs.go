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
	"time"

	"github.com/arnumina/armen.core/pkg/jw"
	"github.com/jackc/pgx/v4"
)

func jobSQLInsert() string {
	return `INSERT INTO jobs (id, name, plugin, type, origin, priority, key, emails, data,
		status, attempt, wf_id, wf_failed, reference, created_at, run_after, result, next_step)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
}

func jobColumns(job *jw.Job) []interface{} {
	return []interface{}{
		job.ID,
		job.Name,
		job.Plugin,
		job.Type,
		job.Origin,
		job.Priority,
		job.Key,
		job.Emails,
		job.Data,
		job.Status,
		job.Attempt,
		job.WfID,
		job.WfFailed,
		job.Reference,
		job.CreatedAt,
		job.RunAfter,
		job.Result,
		job.NextStep,
	}
}

// InsertJob AFAIRE.
func (b *Backend) InsertJob(job *jw.Job) error {
	_, err := b.pgc.Exec(jobSQLInsert(), jobColumns(job)...)
	return err
}

// MaybeInsertJob AFAIRE.
func (b *Backend) MaybeInsertJob(job *jw.Job) (bool, error) {
	t, err := b.pgc.Begin()
	if err != nil {
		return false, err
	}

	defer t.Rollback()

	var id string
	if err := t.QueryRow(
		`SELECT id
		FROM jobs
		WHERE plugin = $1 AND type = $2 AND key = $3 AND (status = $4 OR status = $5 OR status = $6)
		LIMIT 1
		FOR UPDATE`,
		job.Plugin,
		job.Type,
		job.Key,
		jw.Todo,
		jw.Running,
		jw.Pending,
	).Scan(&id); err == nil || !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}

	_, err = t.Exec(jobSQLInsert(), jobColumns(job)...)
	if err != nil {
		return false, err
	}

	if err := t.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

// NextJob AFAIRE.
func (b *Backend) NextJob() (*jw.Job, error) {
	t, err := b.pgc.Begin()
	if err != nil {
		return nil, err
	}

	defer t.Rollback()

	var job jw.Job
	if err := t.QueryRow(
		`SELECT id, name, plugin, type, origin, priority, key, emails, data, status,
		attempt, wf_id, wf_failed, reference, created_at, run_after, result, next_step
		FROM jobs
		WHERE (status = $1 OR status = $2) AND run_after <= $3
		ORDER BY priority DESC, reference ASC
		LIMIT 1
		FOR UPDATE`,
		jw.Todo,
		jw.Pending,
		time.Now(),
	).Scan(
		&job.ID,
		&job.Name,
		&job.Plugin,
		&job.Type,
		&job.Origin,
		&job.Priority,
		&job.Key,
		&job.Emails,
		&job.Data,
		&job.Status,
		&job.Attempt,
		&job.WfID,
		&job.WfFailed,
		&job.Reference,
		&job.CreatedAt,
		&job.RunAfter,
		&job.Result,
		&job.NextStep,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	_, err = t.Exec(`UPDATE jobs SET status = $1 WHERE id = $2`, jw.Running, job.ID)
	if err != nil {
		return nil, err
	}

	if err := t.Commit(); err != nil {
		return nil, err
	}

	return &job, nil
}

// UpdateJob AFAIRE.
func (b *Backend) UpdateJob(job *jw.Job) error {
	_, err := b.pgc.Exec(
		`UPDATE jobs
		SET data = $1, status = $2, attempt = $3, run_after = $4, result = $5, next_step = $6 WHERE id = $7`,
		job.Data,
		job.Status,
		job.Attempt,
		job.RunAfter,
		job.Result,
		job.NextStep,
		job.ID,
	)

	return err
}

/*
######################################################################################################## @(°_°)@ #######
*/
