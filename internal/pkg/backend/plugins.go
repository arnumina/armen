/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package backend

// PluginConfig AFAIRE.
func (b *Backend) PluginConfig(plugin string) (interface{}, error) {
	var cfg interface{}

	if err := b.pgc.QueryRow(`SELECT config FROM plugins WHERE name = $1`, plugin).Scan(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
