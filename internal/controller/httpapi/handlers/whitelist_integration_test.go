package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	mock_service "github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/adapters/mocks"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWhiteList_AddIP(t *testing.T) {
	// Создаем экземпляр объекта WhiteList
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStore := mock_service.NewMockWhiteListStore(controller)

	cases := []struct {
		name    string
		network entity.IPNetwork
	}{
		{name: "valid ip and mask", network: entity.IPNetwork{
			IP:   "192.168.1.1",
			Mask: "255.255.255.0",
		}},
		{name: "invalid ip", network: entity.IPNetwork{
			IP:   "192.12.256.1",
			Mask: "255.255.255.0",
		}},
	}

	for _, testCase := range cases {
		prefix, err := service.GetPrefix(testCase.network.IP, testCase.network.Mask)
		require.NoError(t, err)
		mockStore.EXPECT().Add(prefix, testCase.network.Mask).Return(nil).MaxTimes(1)
		mockStore.EXPECT().Add(prefix, testCase.network.Mask).Return(errors.New("this ip network already exist")).AnyTimes()
	}

	whiteListService := service.NewWhiteList(mockStore, logger)
	whitelist := NewWhiteList(whiteListService, logger)

	// Создаем тестовый HTTP-сервер
	router := httprouter.New()
	router.POST("/auth/whitelist", whitelist.AddIP)

	// Создаем тестовый IP для добавления в белый список
	ip := cases[0].network

	// Кодируем IP в формат JSON
	body, err := json.Marshal(ip)
	require.NoError(t, err)

	ctx := context.Background()
	// Отправляем POST-запрос на сервер для добавления IP в белый список
	req, err := http.NewRequestWithContext(ctx, "POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 204 No Content
	require.Equal(t, http.StatusNoContent, rr.Code)

	// Попытка добавления уже существующего IP в белый список
	req, err = http.NewRequestWithContext(ctx, "POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Попытка добавления невалидного IP в черный список
	invalidIP := cases[1].network
	body, err = json.Marshal(invalidIP)
	require.NoError(t, err)

	req, err = http.NewRequestWithContext(ctx, "POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 400 Bad Request
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWhiteList_RemoveIP(t *testing.T) {
	// Создаем экземпляр объекта WhiteList
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStore := mock_service.NewMockWhiteListStore(controller)

	cases := []struct {
		name    string
		network entity.IPNetwork
	}{
		{name: "valid ip and mask", network: entity.IPNetwork{
			IP:   "192.168.1.1",
			Mask: "255.255.255.0",
		}},
		{name: "invalid ip", network: entity.IPNetwork{
			IP:   "192.12.256.1",
			Mask: "255.255.255.0",
		}},
	}

	for _, testCase := range cases {
		prefix, err := service.GetPrefix(testCase.network.IP, testCase.network.Mask)
		require.NoError(t, err)
		mockStore.EXPECT().Remove(prefix, testCase.network.Mask).Return(nil).AnyTimes()
	}

	whiteListService := service.NewWhiteList(mockStore, logger)
	whitelist := NewWhiteList(whiteListService, logger)

	// Создаем тестовый HTTP-сервер
	router := httprouter.New()
	router.DELETE("/auth/whitelist", whitelist.RemoveIP)

	// Создаем тестовый IP для удаления из белого списка
	ip := cases[0].network

	// Кодируем IP в формат JSON
	body, err := json.Marshal(ip)
	require.NoError(t, err)

	ctx := context.Background()
	// Отправляем POST-запрос на сервер для добавления IP в белый список
	req, err := http.NewRequestWithContext(ctx, "DELETE", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 204 No Content
	require.Equal(t, http.StatusNoContent, rr.Code)

	// Попытка удаления невалидного IP из черного списка
	invalidIP := cases[1].network
	body, err = json.Marshal(invalidIP)
	require.NoError(t, err)

	req, err = http.NewRequestWithContext(ctx, "DELETE", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 400 Bad Request
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWhiteList_ShowIPList(t *testing.T) {
	// Создаем экземпляр объекта BlackList
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStore := mock_service.NewMockWhiteListStore(controller)

	cases := []entity.IPNetwork{
		{
			IP:   "192.168.1.1",
			Mask: "255.255.255.0",
		},
		{
			IP:   "192.168.2.1",
			Mask: "255.255.255.0",
		},
	}

	mockStore.EXPECT().Get().Return(cases, nil).AnyTimes()

	whiteListService := service.NewWhiteList(mockStore, logger)
	whitelist := NewWhiteList(whiteListService, logger)

	// Создаем тестовый HTTP-сервер
	r := httprouter.New()
	r.GET("/auth/whitelist", whitelist.ShowIPList)
	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.Background() // Создаем контекст

	// Отправляем GET-запрос к тестовому серверу с использованием контекста
	req, err := http.NewRequest("GET", ts.URL+"/auth/whitelist", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Проверяем содержимое ответа
	var ipList []entity.IPNetwork
	err = jsoniter.NewDecoder(res.Body).Decode(&ipList)
	require.NoError(t, err)
	assert.Equal(t, cases, ipList)
}
