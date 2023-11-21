package adapters

import "github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/client"

// BlackListRepository репозиторий для работы с черным списком IP-адресов.
type BlackListRepository struct {
	Repository
}

func NewBlackListRepository(client *client.PostgresSQL) *BlackListRepository {
	return &BlackListRepository{Repository: *NewRepository(client, "blacklist")}
}
