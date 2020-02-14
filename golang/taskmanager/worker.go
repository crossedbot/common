package taskmanager

import (
	"fmt"
	"time"

	"github.com/crossedbot/common/golang/logger"
)

type WorkQueue chan Request
type WorkerQueue chan chan Request

type Request struct {
	Task  Task
	Timer *time.Timer
}

func NewRequest(task Task) Request {
	delay := time.Duration(0)
	now := time.Now()
	if now.Before(task.StartTime()) {
		delay = task.StartTime().Sub(now)
	}
	return Request{Task: task, Timer: time.NewTimer(delay)}
}

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	ID          int
	WorkQueue   WorkQueue
	WorkerQueue WorkerQueue
	Quit        chan struct{}
}

func NewWorker(id int, workers WorkerQueue) Worker {
	return &worker{
		ID:          id,
		WorkQueue:   make(WorkQueue),
		WorkerQueue: workers,
		Quit:        make(chan struct{}),
	}
}

func (w *worker) Start() {
	go w.process()
}

func (w *worker) Stop() {
	w.Quit <- struct{}{}
}

func (w *worker) process() {
	logger.Debug(fmt.Sprintf("taskmanager: worker %d started", w.ID))
	for {
		w.WorkerQueue <- w.WorkQueue
		select {
		case work := <-w.WorkQueue:
			go do(work)
		case <-w.Quit:
			logger.Debug(fmt.Sprintf("taskmanager: worker %d stopped", w.ID))
			return
		}
	}
}

func do(r Request) {
	<-r.Timer.C
	r.Task.Call()
}
