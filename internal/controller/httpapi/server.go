package httpapi

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"go.uber.org/zap"
)

// Server представляет HTTP-сервер вашего приложения.
type Server struct {
	server  *http.Server
	handler http.Handler
	config  *config.Config
	logger  *zap.SugaredLogger
}

// NewHTTPAPIServer создает новый экземпляр HttpApiServer с заданными обработчиком, конфигурацией и журналом.
func NewHTTPAPIServer(handler http.Handler, config *config.Config, logger *zap.SugaredLogger) *Server {
	return &Server{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

// Start запускает HTTP-сервер с настройками из конфигурации.
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         s.config.Listen.BindIP + ":" + s.config.Listen.Port,
		Handler:      s.handler,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}
	s.logger.Info("Start http server")
	err := s.server.ListenAndServe()
	return err
}

// ShutdownService завершает работу HTTP-сервера по сигналу операционной системы.
func (s *Server) ShutdownService(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	s.logger.Info("Service is stopping, got signal:", zap.String("signal", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Info(err)
		return
	}
}
