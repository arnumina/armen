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
	"time"
)

func (c *Config) printVersion() error {
	fmt.Println()
	fmt.Println(" ", c.app.Name())
	fmt.Println("===============================================")
	fmt.Println("  version  :", c.app.Version())
	fmt.Println("  built at :", c.app.BuiltAt().String())
	fmt.Println("  by       : Archivage Numérique © INA", time.Now().Year())
	fmt.Println("-----------------------------------------------")
	fmt.Println()

	return ErrStopApp
}

func (c *Config) decryptString() error {
	text, err := c.util.DecryptString(c.decrypt)
	if err != nil {
		return err
	}

	fmt.Println("==>>", text)

	return ErrStopApp
}

func (c *Config) encryptString() error {
	text, err := c.util.EncryptString(c.encrypt)
	if err != nil {
		return err
	}

	fmt.Println("==>>", text)

	return ErrStopApp
}

func (c *Config) handleFlag() error {
	if c.version {
		return c.printVersion()
	}

	if c.key != "" {
		if err := c.util.SetKey(c.key); err != nil {
			return err
		}
	}

	if c.decrypt != "" {
		return c.decryptString()
	}

	if c.encrypt != "" {
		return c.encryptString()
	}

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
