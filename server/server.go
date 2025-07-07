package server

import (
	"os"
	"os/signal"

	"go.uber.org/zap"
)

type Server interface {
	Start() error
	Stop() error
}

type ServerManager struct {
	Servers []Server
	Logger  *zap.Logger
}

// NewServerManager creates a new ServerManager with the provided servers and logger.
func NewServerManager(servers []Server, logger *zap.Logger) *ServerManager {
	return &ServerManager{
		Servers: servers,
		Logger:  logger,
	}
}

func (sm *ServerManager) Run() {
	for _, server := range sm.Servers {
		go func(s Server) {
			if err := s.Start(); err != nil {
				sm.Logger.Error("failed to start server", zap.Error(err))
			}
		}(server)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for {
		select {
		case <-signalChan:
			sm.Logger.Info("received shutdown signal, stopping servers")
			for _, server := range sm.Servers {
				if err := server.Stop(); err != nil {
					sm.Logger.Error("failed to stop server", zap.Error(err))
				}
			}
			return
		}
	}
}
