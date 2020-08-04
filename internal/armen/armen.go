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
	sConfig, err := config.New(a.util, a.app).Load()
	if err != nil {
		if errors.Is(err, config.ErrStopApp) { // --help, --version
			return nil
		}

		return a.onError(err)
	}

	// Logger
	//..................................................................................................................
	sLogger, err := logger.Build(a.util, a.app, sConfig)
	if err != nil {
		return a.onError(err)
	}

	a.ctn.SetLogger(sLogger)

	defer func() {
		sLogger.Info( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"===END",
			"uptime", time.Since(a.app.StartedAt()).Round(time.Second).String(),
		)

		sLogger.Close()
	}()

	sLogger.Info( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"===BEGIN",
		"id", a.app.ID(),
		"name", a.app.Name(),
		"version", a.app.Version(),
		"builtAt", a.app.BuiltAt().String(),
		"pid", os.Getpid(),
	)

	// Bus
	//..................................................................................................................
	sBus := bus.New(a.util, sLogger)
	defer sBus.Close()
	a.ctn.SetBus(sBus)

	// Backend
	//..................................................................................................................
	sBackend, err := backend.New(a.util, sLogger).Build(sConfig)
	if err != nil {
		return a.onError(err)
	}

	defer sBackend.Close()

	// Leader
	//..................................................................................................................
	sLeader, err := leader.New(a.app, sLogger, sBackend).Build()
	if err != nil {
		return a.onError(err)
	}

	defer sLeader.Close()
	a.ctn.SetLeader(sLeader)

	// Model
	//..................................................................................................................
	sModel, err := model.New(sLogger, sBus, sBackend).Build()
	if err != nil {
		return a.onError(err)
	}

	defer sModel.Close()
	a.ctn.SetModel(sModel)

	// Server
	//..................................................................................................................
	sServer, err := server.New(a.util, sLogger, sConfig).Start()
	if err != nil {
		return a.onError(err)
	}

	defer sServer.Stop()
	a.ctn.SetServer(sServer)

	// Plugins
	//..................................................................................................................
	sPlugins, err := plugins.New().Load(a.ctn)
	if err != nil {
		return a.onError(err)
	}

	defer sPlugins.Close()

	// Scheduler
	//..................................................................................................................
	sScheduler, err := scheduler.New(a.util, sLogger, sBus, sLeader).Build(sBackend)
	if err != nil {
		return a.onError(err)
	}

	defer sScheduler.Close()

	// (de)register
	//..................................................................................................................
	if err := sBackend.RegisterInstance(a.app, sServer); err != nil {
		return err
	}

	defer func() {
		_ = sBackend.DeregisterInstance(a.app.ID()) // AFINIR: email ?
	}()

	// Workers
	//..................................................................................................................
	sWorkers := workers.New(a.util, sLogger, sBus, sModel, sPlugins).Start(sConfig)
	defer sWorkers.Stop()

	sScheduler.Start()
	a.waitEnd()
	sLogger.Info("Stopping...") //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
	sScheduler.Stop()

	return nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
