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
	"time"
)

// Lock AFAIRE.
func (b *Backend) Lock(name, owner string, duration time.Duration) (bool, error) {
	now := time.Now()
	expiration := now.Add(duration)

	cTag, err := b.pgc.Exec(
		`UPDATE locks SET expiration_datetime = $1, owner = $2
		WHERE name = $3 AND (owner = $4 OR expiration_datetime IS NULL OR expiration_datetime <= $5)`,
		expiration,
		owner,
		name,
		owner,
		now,
	)
	if err != nil {
		return false, err
	}

	return cTag.RowsAffected() == 1, nil
}

// Unlock AFAIRE.
func (b *Backend) Unlock(name, owner string) error {
	_, err := b.pgc.Exec(
		`UPDATE locks SET expiration_datetime = NULL, owner = NULL WHERE name = $1 AND owner = $2`,
		name,
		owner,
	)

	return err
}

/*
######################################################################################################## @(°_°)@ #######
*/
