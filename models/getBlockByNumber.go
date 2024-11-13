package models

type BlockByNumber struct {
	Result BlockNumberData `json:"result"`
}

type BlockNumberData struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Input any    `json:"input"`
}
