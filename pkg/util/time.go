package util

import (
	"time"
	"strconv"
	// "fmt"
)

func Time() int64 {
	return time.Now().Unix()
}

func GetEra(astr string) int {
	amount, _ := strconv.ParseInt(astr, 10, 64)
	var maxSupply  float64
	var checkpoint float64

	maxSupply  = 3500000000
	checkpoint = maxSupply / 2
	a := float64(amount)
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
