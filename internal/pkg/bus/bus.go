/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package bus

import (
	"regexp"
	"sync"

	"github.com/arnumina/armen.core/pkg/message"
	"github.com/arnumina/armen.core/pkg/resources"
	"github.com/arnumina/logger"

	"github.com/arnumina/armen/internal/pkg/util"
)

const (
	_maxChannelCapacity = 10
	_maxConsumer        = 3
)

type (
	// Resource AFAIRE.
	Resource interface {
		resources.Bus
	}

	// Bus AFAIRE.
	Bus struct {
		util        util.Resource
		logger      *logger.Logger
		subscribers map[*regexp.Regexp]func(*message.Message)
		rwMutex     sync.RWMutex
		group       sync.WaitGroup
	}
)

// New AFAIRE.
func New(util util.Resource, logger *logger.Logger) *Bus {
	return &Bus{
		util:        util,
		logger:      logger,
		subscribers: make(map[*regexp.Regexp]func(*message.Message)),
	}
}

func (b *Bus) goConsumer(publisher string, ch <-chan *message.Message) {
	go func() {
		logger := b.util.CloneLogger(b.logger, publisher)

		logger.Info(">>>Bus") //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

		for msg := range ch {
			msg.Publisher = publisher

			logger.Debug( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
				"Publish message",
				"id", msg.ID,
				"topic", msg.Topic,
			)

			b.rwMutex.RLock()

			for re, cb := range b.subscribers {
				if re.MatchString(msg.Topic) {
					cb(msg)
				}
			}

			b.rwMutex.RUnlock()
		}

		logger.Info("<<<Bus") //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

		b.group.Done()
	}()
}

// AddPublisher AFAIRE.
func (b *Bus) AddPublisher(name string, chCapacity, consumer int) chan<- *message.Message {
	if chCapacity < 0 {
		chCapacity = 0
	} else if chCapacity > _maxChannelCapacity {
		chCapacity = _maxChannelCapacity
	}

	ch := make(chan *message.Message, chCapacity)

	if consumer < 1 {
		consumer = 1
	} else if consumer > _maxConsumer {
		consumer = _maxConsumer
	}

	for n := 0; n < consumer; n++ {
		b.group.Add(1)
		b.goConsumer(name, ch)
	}

	return ch
}

// Subscribe AFAIRE.
func (b *Bus) Subscribe(cb func(*message.Message), reList ...string) error {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	for _, re := range reList {
		regExp, err := regexp.Compile(re)
		if err != nil {
			return err
		}

		b.subscribers[regExp] = cb
	}

	return nil
}

// Close AFAIRE.
func (b *Bus) Close() {
	b.group.Wait()
}

/*
######################################################################################################## @(°_°)@ #######
*/
