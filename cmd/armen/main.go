/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package main

import (
	"os"

	"github.com/arnumina/armen/internal/armen"
)

var (
	_version string
	_builtAt string
)

func main() {
	if armen.New(_version, _builtAt).Run() != nil {
		os.Exit(1)
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
