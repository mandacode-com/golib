package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	errors "github.com/mandacode-com/golib/errors"
)

type Server interface {
	Start() error
	Stop() error
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
	for _, server := range sm.Servers {
		go func(s Server) {
			_ = s.Start()
		}(server)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-signalChan:
	}

	for _, server := range sm.Servers {
		if err := server.Stop(); err != nil {
			return errors.NewPublicError(err.Error(), "failed to stop server", "ERR_SERVER_STOP")
		}
	}
	return nil
}
