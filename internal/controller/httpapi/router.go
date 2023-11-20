package httpapi

import (
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// APIRouter представляет маршрутизатор HTTP API для вашего приложения.
type APIRouter struct {
	router    *httprouter.Router
	auth      *handlers.Authorization
	blackList *handlers.BlackList
	whiteList *handlers.WhiteList
	bucket    *handlers.Bucket
	logger    *zap.SugaredLogger
}

// NewRouter создает новый экземпляр ApiRouter, инициализирует маршрутизатор и связывает обработчики запросов.
func NewRouter(
	auth *handlers.Authorization,
	blackList *handlers.BlackList,
	whiteList *handlers.WhiteList,
	bucket *handlers.Bucket,
	logger *zap.SugaredLogger,
) *APIRouter {
	router := httprouter.New()
	return &APIRouter{
		router:    router,
		auth:      auth,
		blackList: blackList,
		whiteList: whiteList,
		bucket:    bucket,
		logger:    logger,
	}
}

// RegisterRoutes регистрирует маршруты и связывает их с соответствующими обработчиками.
func (r *APIRouter) RegisterRoutes() {
	r.router.POST("/auth/check", r.auth.TryAuthorization)
	r.router.DELETE("/auth/bucket", r.bucket.ResetBucket)
	r.router.POST("/auth/blacklist", r.blackList.AddIP)
	r.router.DELETE("/auth/blacklist", r.blackList.RemoveIP)
	r.router.GET("/auth/blacklist", r.blackList.ShowIPList)
	r.router.POST("/auth/whitelist", r.whiteList.AddIP)
	r.router.DELETE("/auth/whitelist", r.whiteList.RemoveIP)
	r.router.GET("/auth/whitelist", r.whiteList.ShowIPList)
}

// GetRouter возвращает маршрутизатор для использования в вашем приложении.
func (r *APIRouter) GetRouter() *httprouter.Router {
	return r.router
}
