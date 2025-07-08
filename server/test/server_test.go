package server_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/mandacode-com/golib/server"
)

type mockServer struct {
	started bool
	stopped bool
	startMu sync.Mutex
	stopMu  sync.Mutex
	startErr error
	stopErr  error
}

func (m *mockServer) Start() error {
	m.startMu.Lock()
	defer m.startMu.Unlock()
	m.started = true
	return m.startErr
}

func (m *mockServer) Stop() error {
	m.stopMu.Lock()
	defer m.stopMu.Unlock()
	m.stopped = true
	return m.stopErr
}

func TestServerManager_Run_Success(t *testing.T) {
	srv := &mockServer{}
	manager := server.NewServerManager([]server.Server{srv})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := manager.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !srv.started {
		t.Error("expected server to be started")
	}
	if !srv.stopped {
		t.Error("expected server to be stopped")
	}
}

func TestServerManager_Run_StopError(t *testing.T) {
	srv := &mockServer{stopErr: errors.New("stop failed")}
	manager := server.NewServerManager([]server.Server{srv})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := manager.Run(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "stop failed" && err.Error() != "failed to stop server" {
		t.Errorf("unexpected error: %v", err)
	}
	if !srv.stopped {
		t.Error("expected server to be stopped even with error")
	}
}

func TestServerManager_Run_MultipleServers(t *testing.T) {
	srv1 := &mockServer{}
	srv2 := &mockServer{}
	manager := server.NewServerManager([]server.Server{srv1, srv2})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := manager.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !srv1.started || !srv2.started {
		t.Error("expected all servers to be started")
	}
	if !srv1.stopped || !srv2.stopped {
		t.Error("expected all servers to be stopped")
	}
}
