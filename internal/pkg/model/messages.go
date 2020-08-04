/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package model

import "github.com/arnumina/armen.core/pkg/message"

func (m *Model) clean() {
	if err := m.backend.Clean(); err != nil {
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"Clean backend error",
			"reason", err.Error(),
		)
	}
}

func (m *Model) msgHandler(msg *message.Message) {
	m.logger.Info( //:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
		"Consume message",
		"id", msg.ID,
		"topic", msg.Topic,
		"publisher", msg.Publisher,
	)

	switch msg.Topic {
	case "clean":
		m.clean()
	default:
		m.logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"The topic of this message is not valid",
			"id", msg.ID,
			"topic", msg.Topic,
			"publisher", msg.Publisher,
		)
	}
}

func (m *Model) subscribe() error {
	return m.bus.Subscribe(
		m.msgHandler,
		"clean",
	)
}

/*
######################################################################################################## @(°_°)@ #######
*/
