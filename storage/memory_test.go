package storage

import (
	"github.com/canghel3/ethParser/models"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	const address = "0x0"
	mem := NewMemoryStorage()

	t.Run("CREATE", func(t *testing.T) {
		err := mem.Create(address)
		if err != nil {
			t.Errorf("failed to create 0x0 %v", err)
			return
		}

		err = mem.Create(address)
		if err == nil {
			t.Error("expected error already exists")
			return
		}
	})

	t.Run("ADD", func(t *testing.T) {
		tx := models.Transaction{
			From:  address,
			To:    "0x1",
			Input: "hello world",
		}

		err := mem.Add(address, tx)
		if err != nil {
			t.Errorf("failed to add transaction %v", err)
			return
		}

		tx.To = address
		tx.From = "0x2"
		tx.Input = "word up"

		err = mem.Add(address, tx)
		if err != nil {
			t.Errorf("failed to add transaction %v", err)
			return
		}
	})

	t.Run("READ", func(t *testing.T) {
		t.Run("ADDRESS", func(t *testing.T) {
			read, err := mem.Read(address)
			if err != nil {
				t.Errorf("failed to read %v", err)
				return
			}

			if len(read) != 2 {
				t.Errorf("read length should be 2")
				return
			}

			doesNotExist, _ := mem.Read("0x1")
			if len(doesNotExist) != 0 {
				t.Errorf("read length should be 0")
				return
			}
		})

		t.Run("SUBCRIBERS", func(t *testing.T) {
			if len(mem.registry) != 1 {
				t.Errorf("read length should be 2")
				return
			}

			t.Run("TRANSACTIONS", func(t *testing.T) {
				read, _ := mem.Read(address)
				if len(read) != 2 {
					t.Errorf("transactions length should be 2")
					return
				}
			})
		})
	})

	t.Run("DELETE", func(t *testing.T) {
		err := mem.Delete(address)
		if err != nil {
			t.Errorf("failed to delete %v", err)
			return
		}

		tx, _ := mem.Read(address)
		if len(tx) != 0 {
			t.Errorf("transactions should be empty")
			return
		}
	})
}
