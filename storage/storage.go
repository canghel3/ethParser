package storage

import (
	"github.com/canghel3/ethereumBlockchainParser/models"
)

// Storage defines the functionality that any type of storage should implement
// for storing address transactions.
type Storage interface {
	Create(address string) error
	Read(address string) ([]models.Transaction, error)
	Add(address string, transaction models.Transaction) error
	Delete(address string) error
	ReadAllTransactions() []models.Transaction
	ReadAllSubscribers() []string
}
