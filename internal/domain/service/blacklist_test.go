package service

import (
	"testing"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	mock_service "github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/adapters/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddIP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mock_service.NewMockBlackListStore(controller)
	logger := zap.NewExample().Sugar()

	network := entity.IPNetwork{
		IP:   "192.168.1.1",
		Mask: "255.255.255.0",
	}
	prefix, _ := GetPrefix(network.IP, network.Mask)

	store.EXPECT().Add(prefix, network.Mask).Return(nil)

	blackList := NewBlackList(store, logger)
	err := blackList.AddIP(network)
	require.NoError(t, err)
}
