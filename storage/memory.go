package storage

import (
	"fmt"
	"github.com/canghel3/ethParser/models"
	"sync"
)

type MemoryStorage struct {
	mx       sync.RWMutex
	registry map[string][]models.Transaction
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		registry: make(map[string][]models.Transaction),
	}
}

func (ms *MemoryStorage) Create(address string) error {
	if _, ok := ms.registry[address]; !ok {
		ms.mx.Lock()
		ms.registry[address] = make([]models.Transaction, 0)
		ms.mx.Unlock()
		return nil
	}

	return fmt.Errorf("address %s already exists", address)
}

func (ms *MemoryStorage) ReadAllTransactions() []models.Transaction {
	ms.mx.RLock()
	defer ms.mx.RUnlock()
	all := make([]models.Transaction, 0)
	for _, v := range ms.registry {
		all = append(all, v...)
	}

	return all
}

func (ms *MemoryStorage) ReadAllSubscribers() []string {
	var all = make([]string, 0)
	ms.mx.RLock()
	for k := range ms.registry {
		all = append(all, k)
	}

	ms.mx.RUnlock()
	return all
}

func (ms *MemoryStorage) Read(address string) ([]models.Transaction, error) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()
	return ms.registry[address], nil
}

func (ms *MemoryStorage) Add(address string, transaction models.Transaction) error {
	ms.mx.RLock()
	if _, ok := ms.registry[address]; !ok {
		ms.registry[address] = []models.Transaction{}
	}
	ms.mx.RUnlock()

	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.registry[address] = append(ms.registry[address], transaction)
	return nil
}

func (ms *MemoryStorage) Delete(address string) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	delete(ms.registry, address)
	return nil
}
