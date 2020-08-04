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
	"errors"
	"flag"
	"fmt"
	"os"
)

// ErrStopApp AFAIRE.
var ErrStopApp = errors.New("stop application requested")

func (c *Config) printUsage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Println()
		fmt.Println(" ", c.app.Name())
		fmt.Println("================================================================================")
		fs.PrintDefaults()
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println()
	}
}

func (c *Config) parseFlag() error {
	fs := flag.NewFlagSet(c.app.Name(), flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = c.printUsage(fs)

	fs.BoolVar(&c.version, "version", false, "print version information and quit")
	fs.StringVar(&c.cfgString, "config", c.cfgString, "the configuration string")
	fs.StringVar(&c.key, "key", "", "the encryption key")
	fs.StringVar(&c.decrypt, "decrypt", "", "decrypt the string and quit")
	fs.StringVar(&c.encrypt, "encrypt", "", "encrypt the string and quit")
	fs.IntVar(&c.port, "port", 0, "the TCP port")

	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ErrStopApp
		}

		return err
	}

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
