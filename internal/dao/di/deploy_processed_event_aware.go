package di

import (
	"casper-dao-middleware/pkg/casper"
)

type DeployProcessedEventAware struct {
	deployProcessedEvent casper.DeployProcessedEvent
}

func (a *DeployProcessedEventAware) SetDeployProcessedEvent(event casper.DeployProcessedEvent) {
	a.deployProcessedEvent = event
}

func (a *DeployProcessedEventAware) GetDeployProcessedEvent() casper.DeployProcessedEvent {
	return a.deployProcessedEvent
}
