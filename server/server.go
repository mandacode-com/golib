package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	errors "github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
)

type Server interface {
	Start(ctx context.Context) error // Start the server, returns error if it fails
	Stop(ctx context.Context) error  // Stop the server gracefully, returns error if it fails
}

type ServerManager struct {
	Servers []Server
}

func NewServerManager(servers []Server) *ServerManager {
	return &ServerManager{
		Servers: servers,
	}
}

func (sm *ServerManager) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(sm.Servers))

	// Start all servers
	for _, server := range sm.Servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			if err := s.Start(ctx); err != nil {
				errCh <- errors.New(err.Error(), "server start failed", errcode.ErrInternalFailure)
			}
		}(server)
	}

	// Handle shutdown signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-signalChan:
	case err := <-errCh:
		return err
	}

	// Graceful shutdown
	for _, server := range sm.Servers {
		if err := server.Stop(ctx); err != nil {
			return errors.New(err.Error(), "failed to stop server", errcode.ErrInternalFailure)
		}
	}

	// Ensure all Start routines complete
	wg.Wait()
	return nil
}
