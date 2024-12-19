package service

import (
	"fmt"
	"hello/pkg/core/storage"
	"hello/pkg/util"
)

func GetStatus(key string) error {
	oracle := GetOracleService()
	si := storage.GetStatusIndexInstance()

	sf := oracle.storage
	ci := oracle.chain

	sf.Touch()
	ci.Touch()
	si.Load()

	err := sf.Cache()
	key = util.FillHash(key)
	fmt.Println("key: ", key)
	cursor, ok := oracle.storage.CachedUniversalIndexes[key]
	
	if !ok {
			return fmt.Errorf("status bundle not found")
	}

	data, err := oracle.storage.ReadUniversalStatus(cursor)
	if err != nil {
		return fmt.Errorf("err while reading universal status: %s", err)
	}

	fmt.Println(data)
	return nil
}
