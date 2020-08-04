/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package backend

import "github.com/arnumina/armen.core/pkg/model"

// AllEvents AFAIRE.
func (b *Backend) AllEvents() ([]*model.Event, error) {
	rows, err := b.pgc.Query(`SELECT name, disabled, after, repeat FROM events`)
	if err != nil {
		return nil, err
	}

	events := []*model.Event{}

	for rows.Next() {
		var e model.Event

		if err := rows.Scan(
			&e.Name,
			&e.Disabled,
			&e.After,
			&e.Repeat,
		); err != nil {
			return nil, err
		}

		events = append(events, &e)
	}

	return events, nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
