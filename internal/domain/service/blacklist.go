//nolint:dupl //for refactoring
package service

import (
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"go.uber.org/zap"
)

// BlackListStore определяет интерфейс для хранения и управления черным списком IP-сетей.
type BlackListStore interface {
	Add(prefix, mask string) error    // Добавить IP-сеть в черный список.
	Remove(prefix, mask string) error // Удалить IP-сеть из черного списка.
	Get() ([]entity.IPNetwork, error) // Получить список IP-сетей из черного списка.
}

// BlackList представляет сервис управления черным списком IP-сетей.
type BlackList struct {
	store  BlackListStore     // Хранилище черного списка.
	logger *zap.SugaredLogger // Логгер для записи информации.
}

func NewBlackList(store BlackListStore, logger *zap.SugaredLogger) *BlackList {
	return &BlackList{store: store, logger: logger}
}

// processIP обрабатывает операции с IP-сетями, используя указанный метод storeFunc.
func (b *BlackList) processIP(network entity.IPNetwork, storeFunc func(string, string) error) error {
	b.logger.Infoln("Get prefix")
	prefix, err := GetPrefix(network.IP, network.Mask)
	if err != nil {
		return err
	}
	return storeFunc(prefix, network.Mask)
}

// AddIP добавляет IP-сеть в черный список.
func (b *BlackList) AddIP(network entity.IPNetwork) error {
	return b.processIP(network, b.store.Add)
}

// RemoveIP удаляет IP-сеть из черного списка.
func (b *BlackList) RemoveIP(network entity.IPNetwork) error {
	return b.processIP(network, b.store.Remove)
}

// GetIPList возвращает список IP-сетей из черного списка.
func (b *BlackList) GetIPList() ([]entity.IPNetwork, error) {
	return b.store.Get()
}
