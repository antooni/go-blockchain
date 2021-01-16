package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

//Difficulty - number of first bytes that must contain 0s
const Difficulty = 16

//ProofOfWork - allows to calculate Hash of a Block
type ProofOfWork struct {
	Block *Block
	//Target is bigInt because it will work like a bit mask
	//And the number can be pretty long because we convert string to an int
	Target *big.Int
}

//NewProof create new ProofOfWork object, set target, include Block
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//bit-shift 1<<(256-target) e.g.240 so to 240th bit
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

//InitData initialises data, adds nonce, this func will be executed many times
//  until the hash is found
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

//ToHex convert number to hexadecimal, and return it as []byte
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//Run - start proof of work algorithm on a ProofOfWork object
func (pow *ProofOfWork) Run() (int, []byte) {
	// temporary variable to compare hash with target
	// it has to big because we convert string to its numerical value
	var intHash big.Int
	// calculated hash will be put in there
	var hash [32]byte

	// start from nonce equal to zero
	nonce := 0

	// execute until you run out of numbers
	for nonce < math.MaxInt64 {
		// initialise data and calculate hash
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		//fmt.Printf("\r%x", hash)

		// save hash to IntHash
		intHash.SetBytes(hash[:])

		fmt.Printf("\r%x", hash)

		//if we found the hash smaller than the target, it means
		// that it contains enough 0's at the beginning
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else { //otherwise try with incremented noce
			nonce++
		}
	} // execute until you find "perfect and lovely" nonce
	fmt.Println()
	return nonce, hash[:]
}

//Validate the calculated hash, it will be extremely easy with given nonce
// vaidation is way easier than calculation = "Proof of Work"
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1

}

// Trick with difficulty

// set target
// e.g. Difficulty = 4
// length = 12
// target = 000000000001 (1)
// target = target << ( length - Difficilty)=8
// target : 000100000000 (256)

// now we have to find hash smaller than target
// so it will include at least four 0's

// 000100000000 : target
// 000011111111 : max hash value

// so by setting difficulty we are asking miners to find number (nonce)
// which if added to block data will produce hash with at least 4 0's
