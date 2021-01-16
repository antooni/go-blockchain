package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "/tmp/badger"
)

//BlockChain is chain of blocks, connected with hash, such that te following block
// contains the previous block's data
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

//BlockChainIterator is a special type used to retrieve data from key-value db
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//InitBlockChain creates an instance of our blockchain
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions("/tmp/badger")
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err

		}

		item, err := txn.Get([]byte("lh"))
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil

		})
		return err

	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain

}

//AddBlock is a method to add new block to the chain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		return err
	})

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

//Iterator takes Blockchain and converts it to BlockchainIterator
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

//Next - will be called on BlochChainIterator to view the blockchain
// it will be done from last block to genesis
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = append([]byte{}, val...)
			return nil
		})
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
