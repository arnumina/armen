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
	"github.com/arnumina/armen.core/pkg/jw"
)

// Clean AFAIRE.
func (b *Backend) Clean() error {
	t, err := b.pgc.Begin()
	if err != nil {
		return err
	}

	defer t.Rollback()

	_, err = t.Exec(
		`DELETE FROM workflows WHERE status = $1 OR status = $2`,
		jw.Failed,
		jw.Succeeded,
	)
	if err != nil {
		return err
	}

	_, err = t.Exec(
		`DELETE FROM jobs WHERE wf_id IS NULL AND (status = $1 OR status = $2)`,
		jw.Failed,
		jw.Succeeded,
	)
	if err != nil {
		return err
	}

	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
