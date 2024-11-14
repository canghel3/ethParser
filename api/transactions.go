package api

import (
	"encoding/json"
	"github.com/canghel3/ethParser/blockchain"
	"net/http"
)

func transactionsHandler(parser blockchain.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		transactions := parser.GetTransactions(address)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"transactions": transactions,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
