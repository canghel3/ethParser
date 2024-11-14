package api

import (
	"encoding/json"
	"github.com/canghel3/ethereumBlockchainParser/blockchain"
	"log"
	"net/http"
)

type Server struct {
	parser blockchain.Parser
}

func NewServer() *Server {
	return &Server{
		parser: blockchain.NewEthereumParser(),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("POST /subscribe", subscribeHandler(s.parser))
	http.HandleFunc("GET /block", blockHandler(s.parser))
	http.HandleFunc("GET /transactions", transactionsHandler(s.parser))
	http.HandleFunc("GET /all", allHandler(s.parser))

	log.Printf("listening on port %s:%s", "localhost", "1234")
	return http.ListenAndServe(":1234", nil)
}

func allHandler(parser blockchain.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := parser.ReadAllTransactions()
		subscribers := parser.ReadAllSubscribers()

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(map[string]any{
			"transactions": all,
			"subscribers":  subscribers,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
