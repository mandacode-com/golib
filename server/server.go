package server

import (
	"errors"
	"os"
	"os/signal"
)

type Server interface {
	Start() error
	Stop() error
}

type ServerManager struct {
	Servers []Server
}

// NewServerManager creates a new ServerManager with the provided servers and logger.
func NewServerManager(servers []Server) *ServerManager {
	return &ServerManager{
		Servers: servers,
	}
}

func (sm *ServerManager) Run() error {
	for _, server := range sm.Servers {
		go func(s Server) error {
			if err := s.Start(); err != nil {
				return errors.Join(err, errors.New("failed to start server"))
			}
			return nil
		}(server)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for {
		select {
		case <-signalChan:
			for _, server := range sm.Servers {
				if err := server.Stop(); err != nil {
					return errors.Join(err, errors.New("failed to stop server"))
				}
			}
			return nil
		}
	}
}

