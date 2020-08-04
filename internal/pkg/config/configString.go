/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package config

import (
	"fmt"
	"os"
	"strings"
)

func (c *Config) defaultConfigString() error {
	if cs, ok := os.LookupEnv(strings.ToUpper(c.app.Name()) + "_CONFIG"); ok {
		c.cfgString = cs
		return nil
	}

	file := fmt.Sprintf("/etc/%s/%s.yaml", c.app.Name(), c.app.Name())

	if ok, err := c.util.FileExist(file); err != nil {
		return err
	} else if ok {
		c.cfgString = "yaml:file=" + file
		return nil
	}

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
