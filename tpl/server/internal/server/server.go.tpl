package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/xpfo-go/logs"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
)

const (
	defaultGraceTimeout = 30 * time.Second

	defaultIdleTimeout  = 180 * time.Second
	defaultReadTimeout  = 60 * time.Second
	defaultWriteTimeout = 60 * time.Second
)

// Server ...
type Server struct {
	addr     string
	server   *http.Server
	stopChan chan struct{}
	stopOnce sync.Once
	config   *config.Config
}

// NewServer ...
func NewServer(cfg *config.Config, newRouter NewRouterFunc) *Server {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logs.Infof("the server addr: %s", addr)

	// parse the timeouts
	readTimeout := defaultReadTimeout
	writeTimeout := defaultWriteTimeout
	idleTimeout := defaultIdleTimeout

	logs.Infof("the server timeout settings: read_timeout=%s, write_timeout=%s, idle_timeout=%s",
		readTimeout, writeTimeout, idleTimeout)

	router := newRouter(cfg)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &Server{
		addr:     addr,
		server:   server,
		stopChan: make(chan struct{}, 1),
		config:   cfg,
	}
}

// Run ...
func (s *Server) Run(ctx context.Context) {

	go func() {
		<-ctx.Done()
		logs.Info("I have to go...")
		logs.Info("Stopping server gracefully")
		s.Stop()
	}()

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.Errorf("server listen failed: %v", err)
			s.notifyStop()
		}
	}()

	s.Wait()
	logs.Info("Shutting down")
}

// Stop ...
func (s *Server) Stop() {
	defer logs.Info("Server stopped")

	// default graceTimeOut is 60 seconds
	graceTimeout := defaultGraceTimeout

	ctx, cancel := context.WithTimeout(context.Background(), graceTimeout)
	defer cancel()
	logs.Infof("Waiting %s seconds before killing connections...", graceTimeout)

	// disable keep-alive connections
	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		logs.Error(err.Error())
		if closeErr := s.server.Close(); closeErr != nil {
			logs.Errorf("server close failed: %v", closeErr)
		}
	}

	s.notifyStop()
}

// Wait blocks until server is shut down.
func (s *Server) Wait() {
	<-s.stopChan
}

func (s *Server) notifyStop() {
	s.stopOnce.Do(func() {
		s.stopChan <- struct{}{}
	})
}
