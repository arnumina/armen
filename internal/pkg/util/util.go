/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/crypto"
	"github.com/arnumina/logger"
	"github.com/arnumina/uuid"
	"github.com/mitchellh/mapstructure"
)

const (
	_maxNameSize = 10
	_maxIDSize   = 8
)

type (
	// Resource AFAIRE.
	Resource interface {
		resources.Util
		SetKey(key string) error
	}

	// Util AFAIRE.
	Util struct {
		*crypto.Crypto
	}
)

// New AFAIRE.
func New() *Util {
	return &Util{
		Crypto: crypto.New(),
	}
}

// LoggerPrefix AFAIRE.
func (u *Util) LoggerPrefix(name, id string) string {
	if len(name) < _maxNameSize {
		name = strings.Repeat(".", _maxNameSize-len(name)) + name
	} else {
		name = name[:_maxNameSize]
	}

	if len(id) < _maxIDSize {
		id = strings.Repeat(".", _maxIDSize-len(id)) + id
	} else {
		id = id[:_maxIDSize]
	}

	return fmt.Sprintf("%s.%s", name, id)
}

// CloneLogger AFAIRE.
func (u *Util) CloneLogger(logger *logger.Logger, prefix string) *logger.Logger {
	return logger.Clone(u.LoggerPrefix(prefix, uuid.New()))
}

// FileExist AFAIRE.
func (u *Util) FileExist(file string) (bool, error) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// UnixToTime AFAIRE.
func (u *Util) UnixToTime(timestamp string) time.Time {
	ts, err := strconv.ParseInt(timestamp, 0, 64)
	if err != nil {
		ts = 0
	}

	return time.Unix(ts, 0).Local()
}

// DecodeData AFAIRE.
func (u *Util) DecodeData(input, output interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		Result:           output,
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

/*
######################################################################################################## @(°_°)@ #######
*/
