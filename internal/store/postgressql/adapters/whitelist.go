package adapters

import "github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/client"

// WhiteListRepository репозиторий для работы с белым списком IP-адресов.
type WhiteListRepository struct {
	Repository
}

func NewWhiteListRepository(client *client.PostgresSQL) *WhiteListRepository {
	return &WhiteListRepository{Repository: *NewRepository(client, "whitelist")}
}
