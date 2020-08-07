/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package model

// PluginConfig AFAIRE.
func (m *Model) PluginConfig(plugin string, config interface{}) error {
	cfg, err := m.backend.PluginConfig(plugin)
	if err != nil {
		return err
	}

	return m.util.DecodeData(cfg, config)
}

/*
######################################################################################################## @(°_°)@ #######
*/
