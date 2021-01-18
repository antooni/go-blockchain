package blockchain

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

//CanUnlock test the input
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

//CanBeUnlocked tests the output
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
