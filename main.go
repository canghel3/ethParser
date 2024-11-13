package main

import (
	"github.com/canghel3/ethereumBlockchainParser/api"
)

func main() {
	server := api.NewServer()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}