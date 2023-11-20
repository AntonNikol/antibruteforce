package grpcapi

import (
	"context"
	"errors"

	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errInvalidInputIP = errors.New("invalid input IP received from client")

// BlackListServer представляет сервер gRPC для обработки операций с черным списком IP-адресов.
type BlackListServer struct {
	blacklistpb.UnimplementedBlackListServiceServer
	service *service.BlackList
	logger  *zap.SugaredLogger
}

// NewBlackListServer создает новый экземпляр сервера черного списка gRPC.
func NewBlackListServer(service *service.BlackList, logger *zap.SugaredLogger) *BlackListServer {
	return &BlackListServer{service: service, logger: logger}
}

// AddIP обрабатывает запрос на добавление IP-адреса в черный список через gRPC.
func (s *BlackListServer) AddIP(
	_ context.Context,
	req *blacklistpb.AddIpRequest,
) (*blacklistpb.AddIpResponse, error) {
	s.logger.Info("Add IP in blacklist by GRPC")
	ipNetwork := entity.IPNetwork{
		IP:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	IsValidate := handlers.ValidateIP(ipNetwork)
	if !IsValidate {
		s.logger.Info("Invalid input IP received from client")
		return nil, errInvalidInputIP
	}

	err := s.service.AddIP(ipNetwork)
	if err != nil {
		s.logger.Infof("Troubles with add ip: %v", err)
		return nil, err
	}

	res := &blacklistpb.AddIpResponse{IsAddIp: true}
	return res, nil
}

// RemoveIP обрабатывает запрос на удаление IP-адреса из черного списка через gRPC.
func (s *BlackListServer) RemoveIP(
	_ context.Context,
	req *blacklistpb.RemoveIPRequest,
) (*blacklistpb.RemoveIPResponse, error) {
	s.logger.Info("Remove IP in blacklist by GRPC")
	ipNetwork := entity.IPNetwork{
		IP:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	IsValidate := handlers.ValidateIP(ipNetwork)
	if !IsValidate {
		s.logger.Info("Invalid input IP received from client")
		return nil, errInvalidInputIP
	}

	err := s.service.RemoveIP(ipNetwork)
	if err != nil {
		s.logger.Infof("Troubles with remove ip: %v", err)
		return nil, err
	}

	res := &blacklistpb.RemoveIPResponse{IsRemoveIp: true}
	return res, nil
}

// GetIPList обрабатывает запрос на получение списка IP-адресов из черного списка через gRPC.
func (s *BlackListServer) GetIPList(
	_ *blacklistpb.GetIpListRequest,
	stream blacklistpb.BlackListService_GetIpListServer,
) error {
	s.logger.Info("Get IP list in blacklist by GRPC")

	ipList, err := s.service.GetIPList()
	if err != nil {
		s.logger.Infof("Troubles with remove ip: %v", err)
		return err
	}

	for _, network := range ipList {
		err := stream.Send(&blacklistpb.GetIpListResponse{IpNetwork: &blacklistpb.IpNetwork{
			Ip:   network.IP,
			Mask: network.Mask,
		}})
		if err != nil {
			s.logger.Infof("Troubles with remove ip: %v", err)
			return status.Errorf(codes.Internal, "unexpected error: %v", err)
		}
	}

	return nil
}
