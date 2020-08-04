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
	"github.com/arnumina/armen.core/pkg/message"
	"github.com/robfig/cron/v3"
)

type (
	event struct {
		name     string
		disabled bool
		repeat   cron.Schedule
		entryID  cron.EntryID
		tools    *tools
	}
)

func (e *event) Run() {
	if !e.disabled && e.tools.leader.Success() {
		e.Emit()
	}

	if e.entryID != 0 {
		e.tools.cron.Remove(e.entryID)
		e.entryID = 0
	}

	if e.repeat != nil {
		e.tools.cron.Schedule(e.repeat, e)
		e.repeat = nil
	}
}

func (e *event) Emit() {
	msg := message.New(e.name, nil)

	e.tools.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"Emit",
		"event", e.name,
		"message", msg.ID,
	)

	e.tools.channel <- msg
}

func (e *event) Enable() {
	e.disabled = false
}

func (e *event) Disable() {
	e.disabled = true
}

/*
######################################################################################################## @(°_°)@ #######
*/
