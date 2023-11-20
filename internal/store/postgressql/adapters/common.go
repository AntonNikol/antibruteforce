package adapters

import (
	"fmt"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/store/postgressql/client"
)

// Общие константы SQL-запросов для работы с списками IP-адресов.
const (
	isIPExist = `SELECT exists(SELECT 1 FROM %s WHERE prefix = $1 AND mask = $2)`
	insertIP  = `INSERT INTO %s (prefix, mask) VALUES ($1, $2)`
	deleteIP  = `DELETE FROM %s WHERE prefix = $1 AND mask = $2`
	getIPList = `SELECT prefix, mask FROM %s`
)

// Repository представляет репозиторий для работы с списком IP-адресов.
type Repository struct {
	client    *client.PostgresSQL
	tableName string
}

func NewRepository(client *client.PostgresSQL, tableName string) *Repository {
	return &Repository{client: client, tableName: tableName}
}

// Add добавляет IP-адрес в список.
func (r *Repository) Add(prefix string, mask string) error {
	var isExist bool

	err := r.client.DB.QueryRow(fmt.Sprintf(isIPExist, r.tableName), prefix, mask).Scan(&isExist)
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("this ip network already exist")
	}

	err = r.client.DB.QueryRow(fmt.Sprintf(insertIP, r.tableName), prefix, mask).Err()
	if err != nil {
		return err
	}

	return nil
}

// Remove удаляет IP-адрес из списка.
func (r *Repository) Remove(prefix string, mask string) error {
	err := r.client.DB.QueryRow(fmt.Sprintf(deleteIP, r.tableName), prefix, mask).Err()
	if err != nil {
		return err
	}

	return nil
}

// Get возвращает список IP-адресов из списка.
func (r *Repository) Get() ([]entity.IPNetwork, error) {
	ipNetworkList := make([]entity.IPNetwork, 0, 5)

	err := r.client.DB.Select(&ipNetworkList, fmt.Sprintf(getIPList, r.tableName))
	if err != nil {
		return nil, err
	}

	return ipNetworkList, nil
}
