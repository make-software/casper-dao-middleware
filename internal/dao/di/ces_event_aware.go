package di

import (
	"casper-dao-middleware/pkg/go-ces-parser"
)

type CESEventAware struct {
	event ces.Event
}

func (s *CESEventAware) SetCESEvent(event ces.Event) {
	s.event = event
}

func (s *CESEventAware) GetCESEvent() ces.Event {
	return s.event
}
