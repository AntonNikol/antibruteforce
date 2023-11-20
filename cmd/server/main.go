package main

import (
	"log"

	"github.com/AntonNikol/anti-bruteforce/internal/app"
	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadAll()

	log.Println("Success init cfg, start app")
	if err != nil {
		log.Fatalf("can't read env config: %v", err)
		return
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugaredLogger := logger.Sugar()
	defer logger.Sync()

	application := app.NewAntiBruteforceApp(sugaredLogger, cfg)
	application.StartAppAPI()
}
