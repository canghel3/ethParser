package storage

import (
	"github.com/canghel3/ethereumBlockchainParser/models"
)

type Storage interface {
	Create(address string) error
	Read(address string) ([]models.Transaction, error)
	Add(address string, transaction models.Transaction) error
	Delete(address string) error
	ReadAllTransactions() []models.Transaction
	ReadAllSubscribers() []string
}
