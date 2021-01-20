package blockchain

import (
	"bytes"

	"github.com/antooni/go-blockchain/wallet"
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
	Value      int
	PubKeyHash []byte
}

//TxInput is required to create an output, we cannot spend money we dont have
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func NewTXOutput(val int, address string) *TxOutput {
	txo := &TxOutput{val, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
