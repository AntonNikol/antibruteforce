package grpcapi

import (
	"context"
	"errors"

	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	"go.uber.org/zap"
)

// BucketServer представляет сервер gRPC для обработки операций с "бакетом".
type BucketServer struct {
	bucketpb.UnimplementedBucketServiceServer
	service *service.Authorization
	logger  *zap.SugaredLogger
}

// NewBucketServer создает новый экземпляр сервера "бакета" gRPC.
func NewBucketServer(service *service.Authorization, logger *zap.SugaredLogger) *BucketServer {
	return &BucketServer{service: service, logger: logger}
}

// ResetBucket обрабатывает запрос на сброс "бакета" через gRPC.
func (s *BucketServer) ResetBucket(
	_ context.Context,
	req *bucketpb.ResetBucketRequest,
) (*bucketpb.ResetBucketResponse, error) {
	s.logger.Infoln("Reset Bucket by GRPC")

	request := entity.Request{
		Login:    req.GetRequest().GetLogin(),
		Password: req.GetRequest().GetPassword(),
		IP:       req.GetRequest().GetIp(),
	}

	// Устанавливаем пароль в "empty", так как он не используется для сброса "бакета".
	request.Password = "empty"

	// Проверка валидности запроса.
	if !handlers.ValidateRequest(request) {
		return nil, errors.New("invalid input request received from client")
	}

	response := &bucketpb.ResetBucketResponse{}
	isLoginReset := s.service.ResetLoginBucket(request.Login)
	if !isLoginReset {
		response.ResetLogin = false
	} else {
		response.ResetLogin = true
	}

	isIPReset := s.service.ResetIPBucket(request.IP)
	if !isIPReset {
		response.ResetIp = false
	} else {
		response.ResetIp = true
	}

	return response, nil
}
