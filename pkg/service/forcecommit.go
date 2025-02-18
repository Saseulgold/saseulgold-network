package service

import (
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
)

func ForceCommit(txs map[string]*SignedTransaction) error {
	oracle := GetOracleService()
	si := storage.GetStatusIndexInstance()

	sf := oracle.storage
	ci := oracle.chain

	sf.Reset()
	sf.Touch()
	ci.Touch()
	si.Load()

	_, err := oracle.Commit(txs)

	if err != nil {
		return err
	}

	return nil
}
