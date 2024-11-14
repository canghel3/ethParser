package blockchain

import (
	"testing"
	"time"
)

func TestEthereumParser(t *testing.T) {
	const address = "0x0"
	ep := NewEthereumParser()

	//give the parser some time to request the current block, even for slow networks
	time.Sleep(time.Second)

	t.Run("GET CURRENT BLOCK", func(t *testing.T) {
		block := ep.GetCurrentBlock()
		if block == 0 {
			t.Errorf("current block is 0")
			return
		}
	})

	t.Run("SUBSCRIBE", func(t *testing.T) {
		ok := ep.Subscribe(address)
		if !ok {
			t.Errorf("failed to subscribe")
			return
		}

		ok = ep.Subscribe(address)
		if ok {
			t.Errorf("should have already been subscribed")
			return
		}
	})

	time.Sleep(time.Second)

	t.Run("GET TRANSACTION", func(t *testing.T) {
		transactions := ep.GetTransactions(address)
		if len(transactions) != 0 {
			t.Errorf("transactions should be empty")
			return
		}
	})
}
