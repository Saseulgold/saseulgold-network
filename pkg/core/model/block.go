package model

/**
    public $transactions;
    public $s_timestamp;
    public $seal;
    public $universal_updates;
    public $local_updates;

    public $previous_blockhash;
    public $blockhash;
    public $validators;
**/

// txhash: tx
type TransactionMap = map[string]SignedTransaction

type MainBlock struct {
	Height       int64
	Transactions TransactionMap
}

func (block MainBlock) GetBlockRoot() {

}

func (block MainBlock) GetTransactionRoot() {

}

func (block MainBlock) GetUpdateRoot() {

}
