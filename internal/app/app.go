package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/clinterface"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	"github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/adapters"
	"github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/client"
	"go.uber.org/zap"
)

type AntiBruteforceApp struct {
	router                  *httpapi.APIRouter
	grpcBlackListServer     *grpcapi.BlackListServer
	grpcWhiteListServer     *grpcapi.WhiteListServer
	grpcBucketServer        *grpcapi.BucketServer
	grpcAuthorizationServer *grpcapi.AuthorizationServer
	cli                     *clinterface.CommandLineInterface
	clientDB                *client.PostgresSQL
	logger                  *zap.SugaredLogger
	cfg                     *config.Config
}

// NewAntiBruteforceApp создает экземпляр приложения AntiBruteforce и инициализирует все необходимые компоненты.
func NewAntiBruteforceApp(logger *zap.SugaredLogger, cfg *config.Config) *AntiBruteforceApp {
	logger.Infoln("Init http router")
	clientDB := client.NewPostgresSQL(logger, cfg)
	err := clientDB.Open()
	if err != nil {
		logger.Fatalf("Troubels with connect to database: %v", err)
	}

	// Инициализация сервисов и серверов
	blackListStore := adapters.NewBlackListRepository(clientDB)
	blackListService := service.NewBlackList(blackListStore, logger)
	blacklist := handlers.NewBlackList(blackListService, logger)
	grpcBlackListServer := grpcapi.NewBlackListServer(blackListService, logger)

	whitelistStore := adapters.NewWhiteListRepository(clientDB)
	whitelistService := service.NewWhiteList(whitelistStore, logger)
	whitelist := handlers.NewWhiteList(whitelistService, logger)
	grpcWhiteListServer := grpcapi.NewWhiteListServer(whitelistService, logger)

	authorizationService := service.NewAuthorization(blackListService, whitelistService, cfg, logger)
	auth := handlers.NewAuthorization(authorizationService, logger)
	bucket := handlers.NewBucket(authorizationService, logger)
	grpcBucketServer := grpcapi.NewBucketServer(authorizationService, logger)
	grpcAuthorizationServer := grpcapi.NewAuthorizationServer(authorizationService, logger)

	router := httpapi.NewRouter(auth, blacklist, whitelist, bucket, logger)

	cli := clinterface.New(authorizationService, whitelistService, blackListService)
	return &AntiBruteforceApp{
		grpcBlackListServer:     grpcBlackListServer,
		grpcWhiteListServer:     grpcWhiteListServer,
		grpcBucketServer:        grpcBucketServer,
		grpcAuthorizationServer: grpcAuthorizationServer,
		cli:                     cli,
		router:                  router,
		clientDB:                clientDB,
		logger:                  logger,
		cfg:                     cfg,
	}
}

// StartAppAPI запускает приложение AntiBruteforce и выбирает тип сервера (gRPC или HTTP) на основе конфигурации.
func (a *AntiBruteforceApp) StartAppAPI() {
	c := make(chan os.Signal, 1)
	go a.cli.Run(c)
	switch a.cfg.Server.ServerType {
	case "grpc":
		a.logger.Infoln("Init grpc server")
		grpcServer := grpcapi.NewServer(
			a.grpcBlackListServer,
			a.grpcWhiteListServer,
			a.grpcBucketServer,
			a.grpcAuthorizationServer,
			a.cfg,
			a.logger)
		go grpcServer.Shutdown(c)
		err := grpcServer.Start()
		if err != nil {
			a.logger.Fatal(err)
		}
		err = a.clientDB.Close()
		if err != nil {
			fmt.Println(err)
		}

	case "http":
		a.logger.Infoln("Register routes")
		a.router.RegisterRoutes()
		a.logger.Infoln("Init http server")

		server := httpapi.NewHTTPAPIServer(a.router.GetRouter(), a.cfg, a.logger)
		go server.ShutdownService(c)
		err := server.Start()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				a.logger.Infoln(err)
				err = a.clientDB.Close()
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			a.logger.Fatal(err)
		}
	}
}
