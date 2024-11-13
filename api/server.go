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
	http.HandleFunc("POST /subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "address is required", http.StatusBadRequest)
			return
		}

		subcribed := s.parser.Subscribe(address)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]bool{
			"success": subcribed,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("GET /block", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(map[string]int{
			"block": s.parser.GetCurrentBlock(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		transactions := s.parser.GetTransactions(address)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"transactions": transactions,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("GET /all", func(w http.ResponseWriter, r *http.Request) {
		all := s.parser.ReadAll()
		subscribers := s.parser.ReadAllSubscribers()

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(map[string]any{
			"transactions": all,
			"subscribers":  subscribers,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("listening on port %s:%s", "localhost", "1234")
	return http.ListenAndServe(":1234", nil)
}
