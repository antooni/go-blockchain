package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

//Block is a basic bulding block of our blockchain, can cointain variety of data
//the most important is Hash, which will store the result of PoW function
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

//CreateBlock is a wrapper to create a Block
// it contains PoW algorithm inside
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//Genesis the first block is special ^^
func Genesis() *Block {
	return CreateBlock("It is easier to fool people than to convince them that they have been fooled", []byte{})
}

//Serialize converts block data to bytes, do it can be stored in key-value DB (BadgerDB)
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

//Deserialize converts bytes to more "consumable" data type Block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

//Handle is a helper to process errors
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
