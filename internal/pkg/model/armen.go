/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package model

import "github.com/arnumina/armen.core/pkg/model"

// AllInstances AFAIRE.
func (m *Model) AllInstances() ([]*model.Instance, error) {
	return m.backend.AllInstances()
}

/*
######################################################################################################## @(°_°)@ #######
*/
