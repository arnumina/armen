/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package scheduler

import (
	"fmt"
	"time"

	"github.com/arnumina/armen.core/pkg/message"
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"
	"github.com/robfig/cron/v3"

	"github.com/arnumina/armen/internal/pkg/backend"
	"github.com/arnumina/armen/internal/pkg/bus"
	"github.com/arnumina/armen/internal/pkg/leader"
	"github.com/arnumina/armen/internal/pkg/util"
)

type (
	tools struct {
		logger  *logger.Logger
		leader  leader.Resource
		channel chan<- *message.Message
		cron    *cron.Cron
	}

	// Scheduler AFAIRE.
	Scheduler struct {
		tools  *tools
		events map[string]*event
	}
)

// New AFAIRE.
func New(util util.Resource, logger *logger.Logger, bus bus.Resource, leader leader.Resource) *Scheduler {
	return &Scheduler{
		tools: &tools{
			logger:  util.CloneLogger(logger, "scheduler"),
			leader:  leader,
			channel: bus.AddPublisher("scheduler", 1, 1),
		},
	}
}

func (s *Scheduler) build(backend backend.Resource) error {
	events, err := backend.AllEvents()
	if err != nil {
		return err
	}

	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	s.tools.cron = cron.New(cron.WithParser(parser))

	s.events = make(map[string]*event)

	for _, e := range events {
		if e.Name == "" {
			return failure.New(nil).
				Msg("an event name cannot be empty") ///////////////////////////////////////////////////////////////////
		}

		if (e.After == nil || *e.After == "") && (e.Repeat == nil || *e.Repeat == "") {
			return failure.New(nil).
				Set("event", e.Name).
				Msg("this event is not valid") /////////////////////////////////////////////////////////////////////////
		}

		event := &event{
			name:     e.Name,
			disabled: e.Disabled,
			tools:    s.tools,
		}

		s.events[event.name] = event

		if e.After != nil && *e.After != "" {
			d, err := time.ParseDuration(*e.After)
			if err != nil {
				return err
			}

			eID, err := s.tools.cron.AddJob(fmt.Sprintf("@every %s", d.String()), event)
			if err != nil {
				return err
			}

			event.entryID = eID
		}

		if e.Repeat == nil || *e.Repeat == "" {
			continue
		}

		repeat, err := parser.Parse(*e.Repeat)
		if err != nil {
			return err
		}

		if e.After == nil || *e.After == "" {
			_ = s.tools.cron.Schedule(repeat, event)
		} else {
			event.repeat = repeat
		}
	}

	return nil
}

// Build AFAIRE.
func (s *Scheduler) Build(backend backend.Resource) (*Scheduler, error) {
	if err := s.build(backend); err != nil {
		s.Close()
		return nil,
			failure.New(err).Msg("scheduler") //////////////////////////////////////////////////////////////////////////
	}

	return s, nil
}

// Start AFAIRE.
func (s *Scheduler) Start() {
	s.tools.cron.Start()
}

// Stop AFAIRE.
func (s *Scheduler) Stop() {
	s.tools.cron.Stop()
}

// Close AFAIRE.
func (s *Scheduler) Close() {
	close(s.tools.channel)
}

/*
######################################################################################################## @(°_°)@ #######
*/
