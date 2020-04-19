package service

import (
	"context"
	"os"
	"os/signal"
)

type Next func(ctx context.Context) error

type Service interface {
	Run(next Next, sig ...os.Signal) error
	Stop()
}

type service struct {
	ctx    context.Context
	cancel context.CancelFunc
	ch     chan os.Signal
}

func New(ctx context.Context) Service {
	ctx, cancel := context.WithCancel(ctx)
	return &service{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan os.Signal, 1),
	}
}

func (svc *service) Run(next Next, sig ...os.Signal) error {
	if sig != nil && len(sig) > 0 {
		signal.Notify(svc.ch, sig...)
	}
	defer func() {
		signal.Stop(svc.ch)
		svc.cancel()
	}()
	go func() {
		select {
		case <-svc.ch:
			svc.cancel()
		case <-svc.ctx.Done():
		}
		return
	}()
	return next(svc.ctx)
}

func (svc *service) Stop() {
	svc.cancel()
}
