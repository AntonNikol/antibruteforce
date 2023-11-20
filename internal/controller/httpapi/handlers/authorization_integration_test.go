package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	mock_service "github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/adapters/mocks"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTryAuthorization(t *testing.T) {
	// Создаем logger
	logger := zap.NewExample().Sugar()

	// Создаем mock-контроллер
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Создаем mock-объекты для BlackList и WhiteList
	blackListMockStore := mock_service.NewMockBlackListStore(controller)
	whiteListMockStore := mock_service.NewMockWhiteListStore(controller)

	// Создаем конфигурацию
	cfg, err := config.LoadAll()
	require.NoError(t, err)

	// Создаем сервис Authorization с mock-зависимостями
	serviceAuth := service.NewAuthorization(
		service.NewBlackList(blackListMockStore, logger),
		service.NewWhiteList(whiteListMockStore, logger),
		cfg,
		logger,
	)

	// Создаем обработчик
	authorization := NewAuthorization(serviceAuth, logger)

	// Подготовка данных для тестовых случаев
	cases := []struct {
		name    string
		request entity.Request
	}{
		{
			name: "valid request",
			request: entity.Request{
				Login:    "test",
				Password: "1234",
				IP:       "192.1.5.1",
			},
		},
	}

	// Ожидания для моков
	blackListMockStore.EXPECT().Get().Return([]entity.IPNetwork{}, nil).AnyTimes()
	whiteListMockStore.EXPECT().Get().Return([]entity.IPNetwork{}, nil).AnyTimes()

	// Инициализация маршрутизатора и добавление обработчика
	router := httprouter.New()
	router.POST("/auth/check", authorization.TryAuthorization)

	// Основной цикл тестовых случаев
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// Создание HTTP-запроса с JSON-телом
			requestBody, err := json.Marshal(testCase.request)
			require.NoError(t, err)
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, "POST", "/auth/check", bytes.NewReader(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Выполнение HTTP-запроса
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Проверка кода ответа
			require.Equal(t, http.StatusOK, rr.Code)

			// Проверка содержимого ответа
			expectedResponse := "ok=true"
			require.Equal(t, expectedResponse, rr.Body.String())
		})
	}
}
