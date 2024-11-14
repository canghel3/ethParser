package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/canghel3/ethParser/models"
	"github.com/canghel3/ethParser/storage"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ethereumRPCUrl = "https://ethereum-rpc.publicnode.com"

type EthereumParser struct {
	rw              sync.RWMutex
	client          http.Client
	currentBlock    int
	currentBlockHex string
	registry        storage.Storage
}

func NewEthereumParser() *EthereumParser {
	ep := &EthereumParser{
		rw:       sync.RWMutex{},
		client:   http.Client{},
		registry: storage.NewMemoryStorage(),
	}

	go ep.updateBlockNumber()
	return ep
}

func (ep *EthereumParser) ReadAllTransactions() []models.Transaction {
	return ep.registry.ReadAllTransactions()
}

func (ep *EthereumParser) ReadAllSubscribers() []string {
	return ep.registry.ReadAllSubscribers()
}

func (ep *EthereumParser) GetCurrentBlock() int {
	ep.rw.RLock()
	defer ep.rw.RUnlock()
	return ep.currentBlock
}

func (ep *EthereumParser) GetTransactions(address string) []models.Transaction {
	transactions, err := ep.registry.Read(address)
	if err != nil {
		return nil
	}

	return transactions
}

func (ep *EthereumParser) Subscribe(address string) bool {
	if ep.registry.Create(address) == nil {
		go ep.subscriber(address)
		return true
	}

	return false
}

// subscriber is a background process that searches the latest block transactions and
// saves each transaction matching the given address as the sender or receiver of it.
func (ep *EthereumParser) subscriber(address string) {
	var previous string

	defer func() {
		//in case the function panics, restart it and print the panic
		if r := recover(); r != nil {
			go ep.subscriber(address)
			log.Printf("subcriber %s panic: %v", address, r)
		}
	}()

	for {
		if ep.currentBlockHex != previous {
			previous = ep.currentBlockHex
			blockByNumber := ep.getBlockByNumber(ep.currentBlockHex)
			if blockByNumber != nil {
				ep.processBlockTransactions(address, blockByNumber.Result.Transactions)
			}
		}

		time.Sleep(time.Second) //avoid busy-waiting
	}
}

// getBlockByNumber requests specific block information using its hex identifier
func (ep *EthereumParser) getBlockByNumber(hex string) *models.BlockByNumber {
	content := fmt.Sprintf(`{"jsonrpc": "2.0","method":"eth_getBlockByNumber","params":["%s",true],"id":0}`, hex)
	request, err := http.NewRequest("POST", ethereumRPCUrl, strings.NewReader(content))
	if err != nil {
		log.Printf("getBlockByNumber err:%v", err)
		return nil
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := ep.client.Do(request)
	if err != nil {
		log.Printf("getBlockByNumber err:%v", err)
		return nil
	}
	defer response.Body.Close()

	var result models.BlockByNumber
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Printf("getBlockByNumber err:%v", err)
		return nil
	}

	return &result
}

func (ep *EthereumParser) processBlockTransactions(address string, transactions []models.Transaction) {
	for _, transaction := range transactions {
		//process sending and receiving separately
		if transaction.To == address {
			err := ep.registry.Add(address, transaction)
			if err != nil {
				log.Print(err)
				continue
			}
		}

		if transaction.From == address {
			err := ep.registry.Add(address, transaction)
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
}

// updateBlockNumber periodically sends JSONRPC request to update the latest block
func (ep *EthereumParser) updateBlockNumber() {
	//The block time in Ethereum is designed to be approximately 15 seconds. However, it is important to note that block time can vary slightly due to factors such as network congestion and mining difficulty adjustments.
	refreshRate := 5 * time.Second
	sleepTime := 5 * time.Second

	for {
		func() {
			defer func() {
				//in case the function panics, restart it and print the panic
				if r := recover(); r != nil {
					go ep.updateBlockNumber()
					log.Printf("updateBlockNumber panic: %v", r)
				}
			}()

			content := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":0}`
			request, err := http.NewRequest(http.MethodPost, ethereumRPCUrl, strings.NewReader(content))
			if err != nil {
				log.Printf("updateBlockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}

			response, err := ep.client.Do(request)
			if err != nil {
				log.Printf("updateBlockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}
			defer response.Body.Close()

			var blockNumber models.BlockNumber
			err = json.NewDecoder(response.Body).Decode(&blockNumber)
			if err != nil {
				log.Printf("updateBlockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}

			parsedInt, err := strconv.ParseInt(blockNumber.Result[2:], 16, 64)
			if err != nil {
				log.Printf("updateBlockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}

			ep.rw.Lock()
			ep.currentBlock = int(parsedInt)
			ep.currentBlockHex = blockNumber.Result
			ep.rw.Unlock()
		}()

		time.Sleep(refreshRate)
	}
}
