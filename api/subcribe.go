package api

import (
	"encoding/json"
	"github.com/canghel3/ethereumBlockchainParser/blockchain"
	"net/http"
)

func subscribeHandler(parser blockchain.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "address is required", http.StatusBadRequest)
			return
		}

		subcribed := parser.Subscribe(address)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]bool{
			"success": subcribed,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
