//nolint:dupl
package handlers

import (
	"net/http"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// WhiteList представляет обработчик для операций в белом списке IP-адресов.
type WhiteList struct {
	service *service.WhiteList
	logger  *zap.SugaredLogger
}

// NewWhiteList создает новый экземпляр обработчика WhiteList.
func NewWhiteList(service *service.WhiteList, logger *zap.SugaredLogger) *WhiteList {
	return &WhiteList{service: service, logger: logger}
}

// AddIP обрабатывает запрос на добавление IP-адреса в белый список.
func (a *WhiteList) AddIP(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Info("Add IP in whitelist by POST /auth/whitelist")
	initHeaders(rw)
	var inputIP entity.IPNetwork
	err := jsoniter.NewDecoder(r.Body).Decode(&inputIP)
	if err != nil {
		a.logger.Infof("Invalid JSON received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	isValidate := ValidateIP(inputIP)
	if !isValidate {
		a.logger.Info("Invalid input IP received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = a.service.AddIP(inputIP)
	if err != nil {
		if err.Error() == errIPAlreadyExist.Error() {
			a.logger.Info(err)
			rw.WriteHeader(http.StatusBadRequest)
			_, err = rw.Write([]byte(err.Error()))
			if err != nil {
				a.logger.Info(err)
				return
			}
			return
		}
		a.logger.Infof("Troubles with adding IP: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// RemoveIP обрабатывает запрос на удаление IP-адреса из белого списка.
func (a *WhiteList) RemoveIP(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Info("Remove IP in whitelist by POST /auth/whitelist")
	initHeaders(rw)
	var inputIP entity.IPNetwork
	err := jsoniter.NewDecoder(r.Body).Decode(&inputIP)
	if err != nil {
		a.logger.Infof("Invalid JSON received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	isValidate := ValidateIP(inputIP)
	if !isValidate {
		a.logger.Info("Invalid input IP received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = a.service.RemoveIP(inputIP)
	if err != nil {
		a.logger.Infof("Troubles with removing IP: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// ShowIPList обрабатывает запрос на отображение списка IP-адресов в белом списке.
func (a *WhiteList) ShowIPList(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	a.logger.Info("Show IP list in whitelist by GET /auth/whitelist")
	initHeaders(rw)
	ipList, err := a.service.GetIPList()
	if err != nil {
		a.logger.Infof("Troubles with showing IP list: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = jsoniter.NewEncoder(rw).Encode(ipList)
	if err != nil {
		a.logger.Infof("Troubles with encoding IP list: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
