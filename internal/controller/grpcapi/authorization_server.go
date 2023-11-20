package grpcapi

import (
	"context"
	"errors"

	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	"go.uber.org/zap"
)

// AuthorizationServer представляет сервер gRPC для обработки авторизации.
type AuthorizationServer struct {
	authorizationpb.UnimplementedAuthorizationServer
	service *service.Authorization
	logger  *zap.SugaredLogger
}

// NewAuthorizationServer создает новый экземпляр сервера авторизации gRPC.
func NewAuthorizationServer(service *service.Authorization, logger *zap.SugaredLogger) *AuthorizationServer {
	return &AuthorizationServer{service: service, logger: logger}
}

// TryAuthorization обрабатывает запрос на авторизацию через gRPC и возвращает результат.
func (s *AuthorizationServer) TryAuthorization(
	_ context.Context,
	req *authorizationpb.AuthorizationRequest,
) (*authorizationpb.AuthorizationResponse, error) {
	s.logger.Infoln("Try Authorization by GRPC")

	request := entity.Request{
		Login:    req.GetRequest().GetLogin(),
		Password: req.GetRequest().GetPassword(),
		IP:       req.GetRequest().GetIp(),
	}

	// Проверка валидности запроса.
	if !handlers.ValidateRequest(request) {
		return nil, errors.New("invalid input request received from client")
	}

	// Попытка авторизации и проверка на допустимость.
	isAllowed, err := s.service.TryAuthorization(request)
	if err != nil {
		s.logger.Infof("Troubles with authorization request, err: %v", err)
		return nil, err
	}

	return &authorizationpb.AuthorizationResponse{IsAllow: isAllowed}, nil
}
