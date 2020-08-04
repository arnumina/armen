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

	"github.com/arnumina/armen/internal/pkg/application"
	"github.com/arnumina/armen/internal/pkg/server"
)

// RegisterInstance AFAIRE.
func (b *Backend) RegisterInstance(app application.Resource, server server.Resource) error {
	_, err := b.pgc.Exec(
		`INSERT INTO armen (id, host, port, started_at) VALUES ($1, $2, $3, $4)`,
		app.ID(),
		app.FQDN(),
		server.Port(),
		time.Now(),
	)

	return err
}

// DeregisterInstance AFAIRE.
func (b *Backend) DeregisterInstance(id string) error {
	_, err := b.pgc.Exec(`DELETE FROM armen WHERE id = $1`, id)
	return err
}

/*
######################################################################################################## @(°_°)@ #######
*/
