package service

import (
	"fmt"
	"hello/pkg/core/model"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"hello/pkg/util"
)

func CreateGenesisTransaction(privateKey string, publicKey string) (*model.SignedTransaction, error) {
	txData := structure.NewOrderedMap()
	txData.Set("type", "Genesis")
	txData.Set("timestamp", util.UTime())

	tx, err := model.FromRawData(txData, privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	tx.Sign(privateKey, publicKey)

	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func CommitGenesis(privateKey string, publicKey string) error {

	lastHeight := storage.LastHeight()
	if lastHeight > 0 {
		return fmt.Errorf("genesis block already exists. height: %d", lastHeight)
	}

	tx, err := CreateGenesisTransaction(privateKey, publicKey)
	if err != nil {
		panic(err)
	}

	ForceCommit(map[string]*SignedTransaction{tx.GetTxHash(): tx})

	return nil
}
