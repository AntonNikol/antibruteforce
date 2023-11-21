package client

import (
	"fmt"

	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт для использования драйвера PostgreSQL.
	"go.uber.org/zap"
)

// PostgresSQL представляет клиент для работы с базой данных PostgreSQL.
type PostgresSQL struct {
	DB     *sqlx.DB           // Объект для взаимодействия с базой данных.
	logger *zap.SugaredLogger // Логгер для записи информационных сообщений.
	config *config.Config     // Конфигурация для подключения к базе данных.
}

func NewPostgresSQL(logger *zap.SugaredLogger, config *config.Config) *PostgresSQL {
	return &PostgresSQL{logger: logger, config: config}
}

// Open открывает соединение с базой данных PostgreSQL.
func (p *PostgresSQL) Open() error {
	dbSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.config.Database.Host,
		p.config.Database.Port,
		p.config.Database.User,
		p.config.Database.Password,
		p.config.Database.DBName,
		p.config.Database.SslMode,
	)

	db, err := sqlx.Open("postgres", dbSourceName)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	p.DB = db
	p.logger.Info("Connection to db successfully")
	return nil
}

// Close закрывает соединение с базой данных PostgreSQL.
func (p *PostgresSQL) Close() error {
	err := p.DB.Close()
	if err != nil {
		return err
	}
	p.logger.Info("Close db successfully")
	return nil
}
