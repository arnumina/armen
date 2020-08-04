/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package logger

import (
	"github.com/arnumina/failure"
	"github.com/arnumina/logger"

	"github.com/arnumina/armen/internal/pkg/application"
	"github.com/arnumina/armen/internal/pkg/config"
	"github.com/arnumina/armen/internal/pkg/util"
)

func formatter(cfg *config.Logger) (logger.Formatter, error) {
	switch cfg.Formatter {
	case "text":
		return logger.NewTextFormatter(), nil
	default:
		return nil,
			failure.New(nil).
				Set("type", cfg.Formatter).
				Msg("this type of logger formatter does not exist") ////////////////////////////////////////////////////
	}
}

func output(app application.Resource, cfg *config.Logger) (logger.Output, error) {
	switch cfg.Output {
	case "stderr":
		return logger.NewStderrOutput(), nil
	case "stdout":
		return logger.NewStdoutOutput(), nil
	case "syslog":
		return logger.NewSyslogOutput(cfg.Syslog.Facility, app.Name())
	default:
		return nil,
			failure.New(nil).
				Set("type", cfg.Output).
				Msg("this type of logger output does not exist") ///////////////////////////////////////////////////////
	}
}

func build(util util.Resource, app application.Resource, config config.Resource) (*logger.Logger, error) {
	cfg := config.Logger()
	prefix := util.LoggerPrefix(app.Name(), app.ID())

	formatter, err := formatter(cfg)
	if err != nil {
		return nil, err
	}

	output, err := output(app, cfg)
	if err != nil {
		return nil, err
	}

	logger := logger.New(prefix, cfg.Level, formatter, output)

	return logger, nil
}

// Build AFAIRE.
func Build(util util.Resource, app application.Resource, config config.Resource) (*logger.Logger, error) {
	logger, err := build(util, app, config)
	if err != nil {
		return nil,
			failure.New(err).Msg("logger") /////////////////////////////////////////////////////////////////////////////
	}

	return logger, nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
