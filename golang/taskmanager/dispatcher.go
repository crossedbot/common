package taskmanager

import ()

type Dispatcher interface {
	Start()
	Stop()
	Dispatch(r Request)
}

type dispatcher struct {
	WorkQueue   WorkQueue
	MaxRequests int
	WorkerQueue WorkerQueue
	WorkerCount int
	Workers     []Worker
	Quit        chan struct{}
}

func NewDispatcher(nworkers, nrequests int) Dispatcher {
	return &dispatcher{
		WorkQueue:   make(WorkQueue, nrequests),
		MaxRequests: nrequests,
		WorkerQueue: make(WorkerQueue, nworkers),
		WorkerCount: nworkers,
		Quit:        make(chan struct{}),
	}
}

func (d *dispatcher) Start() {
	for i := 0; i < d.WorkerCount; i++ {
		worker := NewWorker(i+1, d.WorkerQueue)
		worker.Start()
		d.Workers = append(d.Workers, worker)
	}
	go d.process()
}

func (d *dispatcher) Stop() {
	for _, w := range d.Workers {
		w.Stop()
	}
	d.Quit <- struct{}{}
}

func (d *dispatcher) Dispatch(r Request) {
	d.WorkQueue <- r
}

func (d *dispatcher) process() {
	for {
		select {
		case work := <-d.WorkQueue:
			go d.take(work)
		case <-d.Quit:
			return
		}
	}
}

func (d *dispatcher) take(r Request) {
	worker := <-d.WorkerQueue
	worker <- r
}
