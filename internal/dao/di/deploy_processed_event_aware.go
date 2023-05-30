package di

import (
	"github.com/make-software/casper-go-sdk/sse"
)

type DeployProcessedEventAware struct {
	deployProcessedEvent sse.DeployProcessedEvent
}

func (a *DeployProcessedEventAware) SetDeployProcessedEvent(event sse.DeployProcessedEvent) {
	a.deployProcessedEvent = event
}

func (a *DeployProcessedEventAware) GetDeployProcessedEvent() sse.DeployProcessedEvent {
	return a.deployProcessedEvent
}
