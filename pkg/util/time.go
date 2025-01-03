package util

import (
	"time"
	// "fmt"
)

func Time() int64 {
	return time.Now().Unix()
}

func GetEra(amount float64) int {
	var maxSupply float64
	var checkpoint float64

	maxSupply =  3500000000
	checkpoint = 1750000000
	a := amount
	era := 0

	for a > maxSupply - checkpoint {
		era++
		maxSupply -= checkpoint
		a -= checkpoint
		checkpoint /= 2
	}

	return era
}





/**
func UTime() int64 {
	microtime := time.Now().UnixNano()
	return ((microtime / 1000) / 1000) * 1000
}
*/
