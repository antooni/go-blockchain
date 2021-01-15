package blockchain

//Block is a basic bulding block of our blockchain, can cointain variety of data
//the most important is Hash, which will store the result of PoW function
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

//BlockChain is chain of blocks, connected with hash, such that te following block
// contains the previous block's data
type BlockChain struct {
	Blocks []*Block
}

//AddBlock is a method to add new block to the chain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
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

//InitBlockChain creates an instance of our blockchain
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
