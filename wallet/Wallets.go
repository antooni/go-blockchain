package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	walletFile = "./tmp/wallets.data"
)

//Wallets is an object which will store all user's addresses and keys
type Wallets struct {
	Wallets map[string]*Wallet
}

//CreateWallets is a wrapper, initialises Wallets, loads values from file
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

//AddWallet allows to generate a new address
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address

}

//GetAllAddresses displays all user's addresses
func (ws *Wallets) GetAllAddresses() []string {
	var list []string

	for address := range ws.Wallets {
		list = append(list, address)
	}

	return list
}

//GetWallet is a simple getter for a wallet from wallets
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

//LoadFile allows loading wallets data from file, adds persistance of data
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}

//SaveFile allows to persist created wallets in a file
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic()
	}

}
