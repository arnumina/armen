/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package armen

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arnumina/armen/internal/pkg/application"
	"github.com/arnumina/armen/internal/pkg/backend"
	"github.com/arnumina/armen/internal/pkg/bus"
	"github.com/arnumina/armen/internal/pkg/config"
	"github.com/arnumina/armen/internal/pkg/container"
	"github.com/arnumina/armen/internal/pkg/leader"
	"github.com/arnumina/armen/internal/pkg/logger"
	"github.com/arnumina/armen/internal/pkg/model"
	"github.com/arnumina/armen/internal/pkg/plugins"
	"github.com/arnumina/armen/internal/pkg/scheduler"
	"github.com/arnumina/armen/internal/pkg/server"
	"github.com/arnumina/armen/internal/pkg/util"
	"github.com/arnumina/armen/internal/pkg/workers"
)

type (
	// Armen AFAIRE.
	Armen struct {
		util util.Resource
		app  application.Resource
		ctn  *container.Container
	}
)

// New AFAIRE.
func New(version, builtAt string) *Armen {
	ctn := container.New()
	util := util.New()
	app := application.New("armen", version, util.UnixToTime(builtAt))

	ctn.SetUtil(util)
	ctn.SetApplication(app)

	return &Armen{
		util: util,
		app:  app,
		ctn:  ctn,
	}
}

func (a *Armen) onError(err error) error {
	fmt.Fprintf( ///////////////////////////////////////////////////////////////////////////////////////////////////////
		os.Stderr,
		"Error: application=%s version=%s builtAt=%s >>> %s\n",
		a.app.Name(),
		a.app.Version(),
		a.app.BuiltAt().String(),
		err,
	)

	return err
}

func (a *Armen) waitEnd() {
	sigEnd := make(chan os.Signal, 1)

	signal.Notify(sigEnd, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTERM)

	<-sigEnd

	close(sigEnd)
}

// Run AFAIRE.
func (a *Armen) Run() error {
	// Application
	//..................................................................................................................
	if err := a.app.Build(); err != nil {
		return a.onError(err)
	}

	// Config
	//..................................................................................................................
	rConfig, err := config.New(a.util, a.app).Load()
	if err != nil {
		if errors.Is(err, config.ErrStopApp) { // --help, --version
			return nil
		}

		return a.onError(err)
	}

	// Logger
	//..................................................................................................................
	rLogger, err := logger.Build(a.util, a.app, rConfig)
	if err != nil {
		return a.onError(err)
	}

	a.ctn.SetLogger(rLogger)

	defer func() {
		rLogger.Info( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"===END",
			"uptime", time.Since(a.app.StartedAt()).Round(time.Second).String(),
		)

		rLogger.Close()
	}()

	rLogger.Info( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"===BEGIN",
		"id", a.app.ID(),
		"name", a.app.Name(),
		"version", a.app.Version(),
		"builtAt", a.app.BuiltAt().String(),
		"pid", os.Getpid(),
	)

	// Bus
	//..................................................................................................................
	rBus := bus.New(a.util, rLogger)
	defer rBus.Close()
	a.ctn.SetBus(rBus)

	// Backend
	//..................................................................................................................
	rBackend, err := backend.New(a.util, rLogger).Build(rConfig)
	if err != nil {
		return a.onError(err)
	}

	defer rBackend.Close()

	// Leader
	//..................................................................................................................
	rLeader, err := leader.New(a.app, rLogger, rBackend).Build()
	if err != nil {
		return a.onError(err)
	}

	defer rLeader.Close()
	a.ctn.SetLeader(rLeader)

	// Model
	//..................................................................................................................
	rModel, err := model.New(rLogger, rBus, rBackend).Build()
	if err != nil {
		return a.onError(err)
	}

	defer rModel.Close()
	a.ctn.SetModel(rModel)

	// Server
	//..................................................................................................................
	rServer, err := server.New(a.util, rLogger, rConfig).Start()
	if err != nil {
		return a.onError(err)
	}

	defer rServer.Stop()
	a.ctn.SetServer(rServer)

	// Plugins
	//..................................................................................................................
	rPlugins, err := plugins.New().Load(a.ctn)
	if err != nil {
		return a.onError(err)
	}

	defer rPlugins.Close()

	// Scheduler
	//..................................................................................................................
	rScheduler, err := scheduler.New(a.util, rLogger, rBus, rLeader).Build(rBackend)
	if err != nil {
		return a.onError(err)
	}

	defer rScheduler.Close()

	// (de)register
	//..................................................................................................................
	if err := rBackend.RegisterInstance(a.app, rServer); err != nil {
		return err
	}

	defer func() {
		_ = rBackend.DeregisterInstance(a.app.ID()) // AFINIR: email ?
	}()

	// Workers
	//..................................................................................................................
	rWorkers := workers.New(a.util, rLogger, rBus, rModel, rPlugins).Start(rConfig)
	defer rWorkers.Stop()

	rScheduler.Start()
	a.waitEnd()
	rLogger.Info("Stopping...") //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
	rScheduler.Stop()

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
