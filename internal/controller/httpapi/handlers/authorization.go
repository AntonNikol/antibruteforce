package handlers

import (
	"net/http"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Authorization представляет обработчик HTTP-запросов для аутентификации.
type Authorization struct {
	service *service.Authorization
	logger  *zap.SugaredLogger
}

// NewAuthorization создает новый экземпляр обработчика аутентификации.
func NewAuthorization(service *service.Authorization, logger *zap.SugaredLogger) *Authorization {
	return &Authorization{service: service, logger: logger}
}

// TryAuthorization обрабатывает POST-запросы на попытку аутентификации.
func (a *Authorization) TryAuthorization(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Infoln("Try Authorization by POST /auth/check")
	initHeaders(rw)
	var request entity.Request
	err := jsoniter.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.logger.Infof("Invalid JSON received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	isValidate := ValidateRequest(request)
	if !isValidate {
		a.logger.Info("Invalid input request received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	isAllowed, err := a.service.TryAuthorization(request)
	if err != nil {
		a.logger.Infof("Troubles with authorization request, err: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isAllowed {
		a.logger.Infof("Request %v is not allowed", request)
		rw.WriteHeader(http.StatusOK)
		_, err = rw.Write([]byte("ok=false"))
		if err != nil {
			a.logger.Infof("Troubles with authorization response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte("ok=true"))
	if err != nil {
		a.logger.Infof("Troubles with authorization response, err: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
