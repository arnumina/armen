/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package backend

import (
	"fmt"
	"time"

	"github.com/arnumina/armen.core/pkg/jw"
	"github.com/arnumina/armen.core/pkg/model"
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"
	"github.com/arnumina/pgsql"

	"github.com/arnumina/armen/internal/pkg/config"
	"github.com/arnumina/armen/internal/pkg/util"
)

const (
	_poolMaxConns = 10
)

type (
	// Resource AFAIRE.
	Resource interface {
		Clean() error
		AllEvents() ([]*model.Event, error)
		PluginConfig(plugin string) (interface{}, error)
		Lock(name, owner string, duration time.Duration) (bool, error)
		Unlock(name, owner string) error
		InsertJob(job *jw.Job) error
		MaybeInsertJob(job *jw.Job) (bool, error)
		NextJob() (*jw.Job, error)
		UpdateJob(job *jw.Job) error
		InsertWorkflow(wf *jw.Workflow, job *jw.Job) error
		Workflow(id string) (*jw.Workflow, error)
		UpdateWorkflow(wf *jw.Workflow) error
	}

	// Backend AFAIRE.
	Backend struct {
		util util.Resource
		pgc  *pgsql.Client
	}
)

// New AFAIRE.
func New(util util.Resource, logger *logger.Logger) *Backend {
	return &Backend{
		util: util,
		pgc:  pgsql.NewClient(util.CloneLogger(logger, "backend")),
	}
}

func (b *Backend) build(config config.Resource) error {
	cfg := config.Backend()

	password, err := b.util.DecryptString(cfg.Password)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		cfg.Username,
		password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		_poolMaxConns,
	)

	if err := b.pgc.Connect(uri); err != nil {
		return err
	}

	return nil
}

// Build AFAIRE.
func (b *Backend) Build(config config.Resource) (*Backend, error) {
	if err := b.build(config); err != nil {
		return nil,
			failure.New(err).Msg("backend") ////////////////////////////////////////////////////////////////////////////
	}

	return b, nil
}

// Close AFAIRE.
func (b *Backend) Close() {
	b.pgc.Close()
}

/*
######################################################################################################## @(°_°)@ #######
*/
