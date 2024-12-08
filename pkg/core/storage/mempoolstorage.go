package storage

import (
	"errors"
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	"sort"
	"sync"
	"time"
)

var (
	ERR_TX_TOO_BIG   = errors.New("transaction size exceeds limit")
	ERR_MEMPOOL_FULL = errors.New("mempool is full")
)

// MempoolTx contains metadata for transactions stored in the mempool
type MempoolTx struct {
	Tx     *SignedTransaction
	Time   int64 // Time when transaction was added to mempool
	Height int   // Block height when transaction was added
	Fee    int64 // Transaction fee
	TxSize int   // Transaction size
}

// MempoolStorage manages the state of the mempool
type MempoolStorage struct {
	mu sync.RWMutex

	// Map using transaction hash as key
	pool map[string]*MempoolTx

	// Transaction priority queue (ordered by ts)
	priorityQueue []*MempoolTx
}

var mempoolInstance *MempoolStorage

// GetMempoolInstance returns a singleton instance of MempoolStorage
func GetMempoolInstance() *MempoolStorage {
	if mempoolInstance == nil {
		mempoolInstance = &MempoolStorage{
			pool:          make(map[string]*MempoolTx),
			priorityQueue: make([]*MempoolTx, 0),
		}
	}
	return mempoolInstance
}

// AddTransaction adds a new transaction to the mempool
func (mp *MempoolStorage) AddTransaction(tx *SignedTransaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	txHash, err := tx.GetTxHash()
	if err != nil {
		return err
	}

	// Check if transaction already exists
	if _, exists := mp.pool[txHash]; exists {
		return nil
	}

	// Check transaction size limit
	size, err := tx.GetSize()
	if err != nil {
		return err
	}
	if size > C.TX_SIZE_LIMIT {
		return ERR_TX_TOO_BIG
	}

	// Check if mempool is full
	if len(mp.pool) >= C.BLOCK_TX_COUNT_LIMIT {
		return ERR_MEMPOOL_FULL
	}

	// Create MempoolTx
	mempoolTx := &MempoolTx{
		Tx:     tx,
		Time:   time.Now().UnixNano(),
		Height: LastHeight(),
		TxSize: size,
	}

	// Add to pool
	mp.pool[txHash] = mempoolTx

	// Add to priority queue
	mp.addToPriorityQueue(mempoolTx)

	return nil
}

// RemoveTransaction removes a transaction from the mempool
func (mp *MempoolStorage) RemoveTransaction(txHash string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if tx, exists := mp.pool[txHash]; exists {
		delete(mp.pool, txHash)
		mp.removeFromPriorityQueue(tx)
	}
}

// GetTransaction retrieves a transaction from the mempool by its hash
func (mp *MempoolStorage) GetTransaction(txHash string) *MempoolTx {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.pool[txHash]
}

// Clear empties the mempool
func (mp *MempoolStorage) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.pool = make(map[string]*MempoolTx)
	mp.priorityQueue = make([]*MempoolTx, 0)
}

// GetTransactions returns a list of transactions sorted by priority
func (mp *MempoolStorage) GetTransactions(limit int) []*SignedTransaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	result := make([]*SignedTransaction, 0, limit)
	for i := 0; i < limit && i < len(mp.priorityQueue); i++ {
		result = append(result, mp.priorityQueue[i].Tx)
	}

	return result
}

// Helper methods for priority queue
func (mp *MempoolStorage) addToPriorityQueue(tx *MempoolTx) {
	mp.priorityQueue = append(mp.priorityQueue, tx)
	mp.sortPriorityQueue()
}

func (mp *MempoolStorage) removeFromPriorityQueue(tx *MempoolTx) {
	for i, item := range mp.priorityQueue {
		if item == tx {
			mp.priorityQueue = append(mp.priorityQueue[:i], mp.priorityQueue[i+1:]...)
			break
		}
	}
}

func (mp *MempoolStorage) sortPriorityQueue() {
	sort.Slice(mp.priorityQueue, func(i, j int) bool {
		return mp.priorityQueue[i].Tx.GetTimestamp() > mp.priorityQueue[j].Tx.GetTimestamp()
	})
}

// FormatTransactions returns a formatted string of transactions in the mempool
func (mp *MempoolStorage) FormatTransactions() string {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	var result string
	for i, tx := range mp.priorityQueue {
		txHash, _ := tx.Tx.GetTxHash()
		result += fmt.Sprintf("%d. TxHash: %s, Time: %s, Height: %d, Size: %d bytes\n",
			i+1,
			txHash,
			time.Unix(0, tx.Time).Format(time.RFC3339),
			tx.Height,
			tx.TxSize,
		)
	}

	if result == "" {
		return "No transactions in mempool"
	}
	return result
}
