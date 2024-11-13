package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/canghel3/ethereumBlockchainParser/models"
	"github.com/canghel3/ethereumBlockchainParser/storage"
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
		rw:           sync.RWMutex{},
		client:       http.Client{},
		currentBlock: 0,
		registry:     storage.NewMemoryStorage(),
	}

	go ep.updateBlockNumber()
	return ep
}

func (ep *EthereumParser) ReadAll() []models.Transaction {
	return ep.registry.ReadAll()
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

func (ep *EthereumParser) subscriber(address string) {
	previous := ep.currentBlockHex

	for {
		if ep.currentBlockHex != previous {
			blockByNumber := ep.getBlockByNumber(ep.currentBlockHex)
			if blockByNumber != nil {
				ep.processBlockTransactions(address, blockByNumber.Result.Transactions)
			}
		}

		time.Sleep(time.Second) //avoid busy-waiting
	}
}

func (ep *EthereumParser) getBlockByNumber(hex string) *models.BlockByNumber {
	content := fmt.Sprintf(`{"jsonrpc": "2.0","method":"eth_getBlockByNumber","params":["%s",true],"id":0}`, hex)
	request, err := http.NewRequest("POST", ethereumRPCUrl, strings.NewReader(content))
	if err != nil {
		log.Printf("eth_getBlockByNumber err:%v", err)
		return nil
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := ep.client.Do(request)
	if err != nil {
		log.Printf("eth_getBlockByNumber err:%v", err)
		return nil
	}
	defer response.Body.Close()

	var result models.BlockByNumber
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Printf("eth_getBlockByNumber err:%v", err)
		return nil
	}

	return &result
}

func (ep *EthereumParser) processBlockTransactions(address string, transactions []models.Transaction) {
	for _, transaction := range transactions {
		//process sending and receiving separately for future-proofing
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

func (ep *EthereumParser) updateBlockNumber() {
	refreshRate := 5 * time.Second
	sleepTime := 10 * time.Second
	for {
		func() {
			content := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":0}`
			request, err := http.NewRequest(http.MethodPost, ethereumRPCUrl, strings.NewReader(content))
			if err != nil {
				log.Printf("eth_blockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}

			response, err := ep.client.Do(request)
			if err != nil {
				log.Printf("eth_blockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}
			defer response.Body.Close()

			var blockNumber models.BlockNumber
			err = json.NewDecoder(response.Body).Decode(&blockNumber)
			if err != nil {
				log.Printf("eth_blockNumber err: %s", err)
				time.Sleep(sleepTime)
				return
			}

			parsedInt, err := strconv.ParseInt(blockNumber.Result[2:], 16, 64)
			if err != nil {
				log.Printf("eth_blockNumber err: %s", err)
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
