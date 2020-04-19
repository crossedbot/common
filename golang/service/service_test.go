package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	svc := New(ctx)
	next := func(ctx context.Context) error {
		return nil
	}
	err := svc.Run(next)
	require.Nil(t, err)
}

func TestStop(t *testing.T) {
	ctx := context.Background()
	svc := New(ctx)
	next := func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}
	go func() {
		time.Sleep(1 * time.Second)
		svc.Stop()
	}()
	err := svc.Run(next)
	require.Nil(t, err)
}
