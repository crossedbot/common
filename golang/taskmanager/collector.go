package taskmanager

import ()

type Collector interface {
	Start()
	Stop()
	Collect(t Task)
}

type collector struct {
	Dispatcher Dispatcher
}

func NewCollector(nworkers, nrequests int) Collector {
	return &collector{Dispatcher: NewDispatcher(nworkers, nrequests)}
}

func (c *collector) Start() {
	c.Dispatcher.Start()
}

func (c *collector) Stop() {
	c.Dispatcher.Stop()
}

func (c *collector) Collect(t Task) {
	c.Dispatcher.Dispatch(NewRequest(t))
}
