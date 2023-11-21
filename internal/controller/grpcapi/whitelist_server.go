package grpcapi

import (
	"context"

	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/whitelistpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WhiteListServer представляет сервер gRPC для операций с белым списком IP-адресов.
type WhiteListServer struct {
	whitelistpb.UnimplementedWhiteListServiceServer
	service *service.WhiteList
	logger  *zap.SugaredLogger
}

// NewWhiteListServer создает новый экземпляр сервера gRPC для белого списка IP-адресов.
func NewWhiteListServer(service *service.WhiteList, logger *zap.SugaredLogger) *WhiteListServer {
	return &WhiteListServer{service: service, logger: logger}
}

// AddIP обрабатывает запрос на добавление IP-адреса в белый список через gRPC.
func (s *WhiteListServer) AddIP(
	_ context.Context,
	req *whitelistpb.AddIpRequest,
) (*whitelistpb.AddIpResponse, error) {
	s.logger.Info("Add IP in whitelist by GRPC")
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

	res := &whitelistpb.AddIpResponse{IsAddIp: true}
	return res, nil
}

// RemoveIP обрабатывает запрос на удаление IP-адреса из белого списка через gRPC.
func (s *WhiteListServer) RemoveIP(
	_ context.Context,
	req *whitelistpb.RemoveIPRequest,
) (*whitelistpb.RemoveIPResponse, error) {
	s.logger.Info("Remove IP in whitelist by GRPC")
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

	res := &whitelistpb.RemoveIPResponse{IsRemoveIp: true}
	return res, nil
}

// GetIPList обрабатывает запрос на получение списка IP-адресов из белого списка через gRPC.
func (s *WhiteListServer) GetIPList(
	_ *whitelistpb.GetIpListRequest,
	stream whitelistpb.WhiteListService_GetIpListServer,
) error {
	s.logger.Info("Get IP list in whitelist by GRPC")

	ipList, err := s.service.GetIPList()
	if err != nil {
		s.logger.Infof("Troubles with remove ip: %v", err)
		return err
	}

	for _, network := range ipList {
		err := stream.Send(&whitelistpb.GetIpListResponse{IpNetwork: &whitelistpb.IpNetwork{
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
