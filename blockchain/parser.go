package blockchain

import "github.com/canghel3/ethereumBlockchainParser/models"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []models.Transaction
	ReadAllTransactions() []models.Transaction
	ReadAllSubscribers() []string
}
