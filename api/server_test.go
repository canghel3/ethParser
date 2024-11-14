package api

import (
	"encoding/json"
	"github.com/canghel3/ethereumBlockchainParser/blockchain"
	"github.com/canghel3/ethereumBlockchainParser/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	const address = "0x0"

	ep := blockchain.NewEthereumParser()

	t.Run("SUBSCRIBE", func(t *testing.T) {
		t.Run("DOES NOT EXIST", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/subscribe?address="+address, nil)
			w := httptest.NewRecorder()

			subscribeHandler(ep).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
				return
			}

			var result map[string]bool
			json.NewDecoder(resp.Body).Decode(&result)
			if !result["success"] {
				t.Errorf("expected successful subcription")
				return
			}
		})

		t.Run("ALREADY EXIST", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/subscribe?address="+address, nil)
			w := httptest.NewRecorder()

			subscribeHandler(ep).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
				return
			}

			var result map[string]bool
			json.NewDecoder(resp.Body).Decode(&result)
			if result["success"] {
				t.Error("should have already been subcribed")
				return
			}
		})

	})

	t.Run("BLOCK", func(t *testing.T) {
		//give it some time to request the current block
		time.Sleep(time.Second)
		req := httptest.NewRequest("GET", "/block", nil)
		w := httptest.NewRecorder()

		blockHandler(ep).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
			return
		}

		var result map[string]int
		json.NewDecoder(resp.Body).Decode(&result)
		if result["block"] == 0 {
			t.Error("block should not be 0")
			return
		}
	})

	t.Run("TRANSACTIONS", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/transactions", nil)
		w := httptest.NewRecorder()

		transactionsHandler(ep).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
			return
		}

		var result map[string][]models.Transaction
		json.NewDecoder(resp.Body).Decode(&result)
		if len(result["transactions"]) != 0 {
			t.Error("transactions length should be 0")
			return
		}
	})
}
