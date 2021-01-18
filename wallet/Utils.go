package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode
}

//Base58 was invented with bitcoin, it is similiar to Base64 but...
// does not inlude 0,O,l,I,+,/ in its alphabet (6 less than Base64)
// it was invented mainly to avoid problems with similarily looking letters
