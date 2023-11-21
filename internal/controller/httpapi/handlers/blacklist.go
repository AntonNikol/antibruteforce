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

// BlackList представляет обработчик HTTP-запросов для работы с черным списком IP-адресов.
type BlackList struct {
	service *service.BlackList
	logger  *zap.SugaredLogger
}

// NewBlackList создает новый экземпляр обработчика черного списка.
func NewBlackList(service *service.BlackList, logger *zap.SugaredLogger) *BlackList {
	return &BlackList{service: service, logger: logger}
}

// AddIP обрабатывает POST-запрос на добавление IP-адреса в черный список.
func (a *BlackList) AddIP(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Info("Add IP in blacklist by POST /auth/blacklist")
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
		a.logger.Infof("Troubles with add IP: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// RemoveIP обрабатывает POST-запрос на удаление IP-адреса из черного списка.
func (a *BlackList) RemoveIP(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Info("Remove IP in blacklist by POST /auth/blacklist")
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
		a.logger.Infof("Troubles with remove IP: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// ShowIPList обрабатывает GET-запрос на получение списка IP-адресов из черного списка.
func (a *BlackList) ShowIPList(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	a.logger.Info("Show IP list in blacklist by GET /auth/blacklist")
	initHeaders(rw)
	ipList, err := a.service.GetIPList()
	if err != nil {
		a.logger.Infof("Troubles with show IP list: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = jsoniter.NewEncoder(rw).Encode(ipList)
	if err != nil {
		a.logger.Infof("Troubles with encode IP list: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
