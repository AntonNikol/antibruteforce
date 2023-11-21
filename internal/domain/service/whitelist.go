//nolint:dupl //for refactor
package service

import (
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"go.uber.org/zap"
)

// WhiteListStore представляет интерфейс хранилища белого списка IP-адресов.
type WhiteListStore interface {
	Add(prefix, mask string) error    // Добавляет IP-адрес в белый список.
	Remove(prefix, mask string) error // Удаляет IP-адрес из белого списка.
	Get() ([]entity.IPNetwork, error) // Получает список IP-адресов из белого списка.
}

// WhiteList представляет сервис для управления белым списком IP-адресов.
type WhiteList struct {
	store  WhiteListStore     // Хранилище белого списка IP-адресов.
	logger *zap.SugaredLogger // Логгер.
}

// NewWhiteList создает новый экземпляр WhiteList с указанным хранилищем и логгером.
func NewWhiteList(store WhiteListStore, logger *zap.SugaredLogger) *WhiteList {
	return &WhiteList{store: store, logger: logger}
}

// processIP обрабатывает операции с IP-адресами, используя указанный метод storeFunc.
func (w *WhiteList) processIP(network entity.IPNetwork, storeFunc func(string, string) error) error {
	w.logger.Infoln("Get prefix")
	prefix, err := GetPrefix(network.IP, network.Mask) // Получить префикс IP-адреса.
	if err != nil {
		return err
	}
	return storeFunc(prefix, network.Mask) // Выполнить операцию белого списка, используя префикс и маску.
}

// AddIP добавляет IP-адрес в белый список.
func (w *WhiteList) AddIP(network entity.IPNetwork) error {
	return w.processIP(network, w.store.Add)
}

// RemoveIP удаляет IP-адрес из белого списка.
func (w *WhiteList) RemoveIP(network entity.IPNetwork) error {
	return w.processIP(network, w.store.Remove)
}

// GetIPList возвращает список IP-адресов из белого списка.
func (w *WhiteList) GetIPList() ([]entity.IPNetwork, error) {
	return w.store.Get() // Получить список IP-адресов из хранилища белого списка.
}
