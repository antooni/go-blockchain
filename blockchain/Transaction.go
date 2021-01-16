package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

//Transaction is a basic data structure used inside a block
// it allows validating transfers of coins
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

//TxOutput is like a receipt, how much and to whom
type TxOutput struct {
	Value  int
	PubKey string
}

//TxInput is required to create an output, we cannot spend money we dont have
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

//SetID calculates the hash of a transaction and sets it as an ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//CoinbaseTx is a special method for a miner to generate new coins
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

//IsCoinbase is a helper function to determine wether transaction is a coinbase-tx
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

//CanUnlock test the input
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

//CanBeUnlocked tests the output
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
