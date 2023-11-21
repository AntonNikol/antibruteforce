package handlers

import (
	"net/http"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Bucket представляет обработчик HTTP-запросов для сброса bucket.
type Bucket struct {
	service *service.Authorization
	logger  *zap.SugaredLogger
}

// NewBucket создает новый экземпляр обработчика bucket.
func NewBucket(service *service.Authorization, logger *zap.SugaredLogger) *Bucket {
	return &Bucket{service: service, logger: logger}
}

// ResetBucket обрабатывает POST-запрос для сброса bucket.
func (a *Bucket) ResetBucket(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Infoln("Reset Bucket by POST /auth/bucket")
	initHeaders(rw)
	var request entity.Request
	err := jsoniter.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.logger.Infof("Invalid JSON received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	request.Password = "empty"
	isValidate := ValidateRequest(request)
	if !isValidate {
		a.logger.Info("Invalid input request received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	isLoginReset := a.service.ResetLoginBucket(request.Login)
	if !isLoginReset {
		_, err = rw.Write([]byte("resetLogin=false\n"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		_, err = rw.Write([]byte("resetLogin=true\n"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	isIPReset := a.service.ResetIPBucket(request.IP)
	if !isIPReset {
		_, err = rw.Write([]byte("resetIp=false"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		_, err = rw.Write([]byte("resetIp=true"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
