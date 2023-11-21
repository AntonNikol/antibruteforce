package grpcapi

import (
	"net"
	os "os"
	"os/signal"
	"syscall"

	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/whitelistpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server представляет сервер gRPC для обработки запросов.
type Server struct {
	blackListServer     blacklistpb.BlackListServiceServer
	whiteListServer     whitelistpb.WhiteListServiceServer
	bucketServer        bucketpb.BucketServiceServer
	authorizationServer authorizationpb.AuthorizationServer
	grpcServer          *grpc.Server
	config              *config.Config
	logger              *zap.SugaredLogger
}

// NewServer создает новый экземпляр сервера gRPC.
func NewServer(
	blackListServer blacklistpb.BlackListServiceServer,
	whiteListServer whitelistpb.WhiteListServiceServer,
	bucketServer bucketpb.BucketServiceServer,
	authorizationServer authorizationpb.AuthorizationServer,
	config *config.Config,
	logger *zap.SugaredLogger,
) *Server {
	grpcServer := grpc.NewServer()
	return &Server{
		blackListServer:     blackListServer,
		whiteListServer:     whiteListServer,
		bucketServer:        bucketServer,
		authorizationServer: authorizationServer,
		config:              config,
		grpcServer:          grpcServer,
		logger:              logger,
	}
}

// Start запускает сервер gRPC и начинает прослушивать указанный порт.
func (s *Server) Start() error {
	s.logger.Infoln("start grpc server")
	listener, err := net.Listen("tcp", s.config.Listen.BindIP+":"+s.config.Listen.Port)
	if err != nil {
		return err
	}
	blacklistpb.RegisterBlackListServiceServer(s.grpcServer, s.blackListServer)
	whitelistpb.RegisterWhiteListServiceServer(s.grpcServer, s.whiteListServer)
	bucketpb.RegisterBucketServiceServer(s.grpcServer, s.bucketServer)
	authorizationpb.RegisterAuthorizationServer(s.grpcServer, s.authorizationServer)
	reflection.Register(s.grpcServer)
	err = s.grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

// Shutdown останавливает сервер gRPC при получении сигнала завершения.
func (s *Server) Shutdown(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	s.logger.Info("Service is stop, got signal:", zap.String("signal", sig.String()))
	s.grpcServer.GracefulStop()
}
