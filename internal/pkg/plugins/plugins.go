/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	_plugin "plugin"

	"github.com/arnumina/armen.core/pkg/plugin"
	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/failure"
)

const (
	pluginFuncName = "Export"
)

type (
	// Resource AFAIRE.
	Resource interface {
		Find(name string) plugin.Plugin
	}

	// Plugins AFAIRE.
	Plugins struct {
		all map[string]plugin.Plugin
	}
)

// New AFAIRE.
func New() *Plugins {
	return &Plugins{
		all: make(map[string]plugin.Plugin),
	}
}

func (p *Plugins) loadOne(ctn resources.Container, path string) error {
	plg, err := _plugin.Open(path)
	if err != nil {
		return err
	}

	symbol, err := plg.Lookup(pluginFuncName)
	if err != nil {
		return err
	}

	fn, ok := symbol.(func(resources.Container) plugin.Plugin)
	if !ok {
		return failure.New(nil).
			Set("plugin", path).
			Msg("this plugin doesn't export the right function") ///////////////////////////////////////////////////////
	}

	plugin := fn(ctn)

	_, ok = p.all[plugin.Name()]
	if ok {
		return failure.New(nil).
			Set("plugin", path).
			Msg("another plugin with the same name already exists") ////////////////////////////////////////////////////
	}

	if err := plugin.Build(); err != nil {
		return failure.New(err).
			Set("plugin", path).
			Msg("impossible to build this plugin") /////////////////////////////////////////////////////////////////////
	}

	p.all[plugin.Name()] = plugin

	ctn.Logger().Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"Plugin",
		"name", plugin.Name(),
		"version", plugin.Version(),
		"builtAt", plugin.BuiltAt().String(),
	)

	return nil
}

func (p *Plugins) load(ctn resources.Container) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	plugins, err := filepath.Glob(filepath.Join(filepath.Dir(exe), fmt.Sprintf("%s.*.so", ctn.Application().Name())))
	if err != nil {
		return err
	}

	for _, path := range plugins {
		if err := p.loadOne(ctn, path); err != nil {
			return err
		}
	}

	return nil
}

// Load AFAIRE.
func (p *Plugins) Load(ctn resources.Container) (*Plugins, error) {
	if err := p.load(ctn); err != nil {
		p.Close()
		return nil,
			failure.New(err).Msg("plugins") ////////////////////////////////////////////////////////////////////////////
	}

	return p, nil
}

// Find AFAIRE.
func (p *Plugins) Find(name string) plugin.Plugin {
	return p.all[name]
}

// Close AFAIRE.
func (p *Plugins) Close() {
	for _, plugin := range p.all {
		plugin.Close()
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
