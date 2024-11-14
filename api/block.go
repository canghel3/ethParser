package api

import (
	"encoding/json"
	"github.com/canghel3/ethereumBlockchainParser/blockchain"
	"net/http"
)

func blockHandler(parser blockchain.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(map[string]int{
			"block": parser.GetCurrentBlock(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
