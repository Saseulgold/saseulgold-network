package service

import (
	"context"
	"encoding/json"
	"fmt"
	C "hello/pkg/core/config"
	"hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"hello/pkg/core/vm"
	"hello/pkg/swift"
	"hello/pkg/util"
	"os"
	"time"

	"sync"

	"go.uber.org/zap"
)

var logger = util.GetLogger()

type OracleState string

const (
	StateCommitting   OracleState = "committing"
	StateTransaction  OracleState = "transaction"
	StateInitializing OracleState = "initializing"
	StateStopped      OracleState = "stopped"
)

func OracleLog(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println()
}

type Oracle struct {
	swift        *swift.Server
	machine      *vm.Machine
	storage      *storage.StatusFile
	storageIndex *storage.StatusIndex
	chain        *storage.ChainStorage
	mempool      *storage.MempoolStorage

	mu sync.RWMutex

	State    OracleState
	Replicas map[string]int64
}

var oracleInstance *Oracle

func GetOracleService() *Oracle {
	if oracleInstance == nil {
		oracleInstance = &Oracle{
			machine:      vm.GetMachineInstance(),
			mempool:      storage.GetMempoolInstance(),
			chain:        storage.GetChainStorageInstance(),
			swift:        nil,
			storage:      storage.GetStatusFileInstance(),
			storageIndex: storage.GetStatusIndexInstance(),
			State:        StateTransaction,
			Replicas:     map[string]int64{},
		}
	}
	return oracleInstance
}

func (o *Oracle) Consensus(txs map[string]*model.SignedTransaction) {
	o.mu.Lock()
	defer o.mu.Unlock()

}

func (o *Oracle) Commit(txs map[string]*model.SignedTransaction) (*model.Block, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	var previousBlockhash string

	lastBlockHeight := storage.LastHeight()
	previousBlock, err := storage.GetChainStorageInstance().GetBlock(lastBlockHeight)

	if previousBlock == nil {
		previousBlockhash = ""
	} else {
		previousBlockhash = previousBlock.BlockHash()
	}

	o.machine.Init(previousBlock)
	o.machine.SetTransactions(txs)
	o.machine.PreCommit()

	block := model.NewBlock(storage.LastHeight()+1, previousBlockhash)
	block.SetTimestamp(int64(util.Utime()))

	expectedBlock := o.machine.ExpectedBlock()

	if expectedBlock.GetTransactionCount() == 0 {
		OracleLog("no transactions to commit. invalid block.")
		return nil, fmt.Errorf("no transactions to commit. invalid block.")
	}

	err = o.machine.Commit(expectedBlock)

	if err != nil {
		OracleLog("Commit error: %v", err)
		return nil, fmt.Errorf("Commit error: %v", err)
	}

	transactions := expectedBlock.GetTransactions()
	txHashes := make([]string, 0, len(transactions))
	for txHash := range transactions {
		txHashes = append(txHashes, txHash)
	}

	logger.Info("Commit success", zap.Int("transaction_count", len(txHashes)))
	return expectedBlock, nil

}

func (o *Oracle) OnStartCommit() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.State = StateCommitting
}

func (o *Oracle) OnFinishCommit() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.State = StateTransaction
}

func (o *Oracle) Run() error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			switch o.State {

			case StateTransaction:
				o.OnStartCommit()
				transactions := o.mempool.PopTransactionsHashMap()

				if len(transactions) == 0 {
					o.OnFinishCommit()
					continue
				}

				fmt.Println("transactions: ", transactions, len(transactions))

				block, err := o.Commit(transactions)
				o.BroadcastBlock(block)

				o.OnFinishCommit()

				if err != nil {
					OracleLog("failed to commit transactions: %v", err)
				}
			}
		}
	}
}

func (o *Oracle) OnStartUp(config swift.SecurityConfig, port string) error {
	/**
	isRunning := util.ServiceIsRunning(storage.DataRootDir(), "oracle")
	if isRunning {
		return fmt.Errorf("oracle is already running")
	}
	**/

	o.chain.Touch()
	o.storage.Touch()

	o.storage.Reset()
	o.storageIndex.Load()

	o.swift = swift.NewServer(C.SWIFT_HOST+":"+port, config)
	o.registerPacketHandlers()
	err := util.ProcessStart(storage.DataRootDir(), "oracle", os.Getpid())
	if err != nil {
		OracleLog("failed to start oracle: %v", err)
		return err
	}

	return nil
}

func (o *Oracle) Start() error {
	go o.Run()
	err := o.swift.Start()
	if err != nil {
		return err
	}

	return nil
}

func (o *Oracle) registerPacketHandlers() {

	// send trasaction a node
	o.swift.RegisterHandler(swift.PacketTypeSendTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		swift.SwiftInfoLog("send transaction request: %s", string(packet.Payload))

		data, err := structure.ParseOrderedMap(string(packet.Payload))
		if err != nil {
			return err
		}

		tx, err := model.NewSignedTransaction(data)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		txHash := tx.GetTxHash()
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		err = o.mempool.AddTransaction(&tx)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		err = o.swift.Send(ctx, packet)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeSendTransactionResponse,
			Payload: json.RawMessage(fmt.Sprintf("%v", txHash)),
		}

		return o.swift.Send(ctx, response)
	})

	// broadcast transaction request
	o.swift.RegisterHandler(swift.PacketTypeBroadcastTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		swift.SwiftInfoLog("broadcast transaction request: %s", string(packet.Payload))

		var tx model.SignedTransaction
		data, err := structure.ParseOrderedMap(string(packet.Payload))

		if err != nil {
			return err
		}

		tx, err = model.NewSignedTransaction(data)

		if err != nil {
			return err
		}

		lastBlockHeight := storage.LastHeight()
		previousBlock, err := storage.GetChainStorageInstance().GetBlock(lastBlockHeight)

		machine := vm.NewMachine(previousBlock)
		valid, err := machine.TxValidity(&tx)

		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		if !valid {
			return o.swift.SendErrorResponse(ctx, "transaction is not valid")
		}

		err = o.mempool.AddTransaction(&tx)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		txHash := tx.GetTxHash()
		logger.Info("successfully added transaction to mempool", zap.String("tx_hash", txHash))

		var responseData []byte
		swift.SwiftInfoLog("broadcasting transaction to peers: %s", o.swift.GetPeers())

		// if err := o.swift.Broadcast(ctx, packet); err != nil {
		// 	return o.swift.SendErrorResponse(ctx, err.Error())
		// }

		responseData, err = json.Marshal(map[string]string{"ok": "true", "msg": "successfully broadcast transaction"})

		if err != nil {
			return err
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeBroadcastTransactionResponse,
			Payload: json.RawMessage(responseData),
		}

		return o.swift.Send(ctx, response)
	})

	// list mempool transaction request
	o.swift.RegisterHandler(swift.PacketTypeListMempoolTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		data, err := json.Marshal(o.mempool.FormatTransactions())
		if err != nil {
			return err
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeListMempoolTransactionResponse,
			Payload: json.RawMessage(data),
		}

		return o.swift.Send(ctx, response)
	})

	o.swift.RegisterHandler(swift.PacketTypeGetStatusBundleRequest, func(ctx context.Context, packet *swift.Packet) error {

		var payload map[string]string
		err := json.Unmarshal(packet.Payload, &payload)
		if err != nil {
			return err
		}

		cursor, ok := o.storage.CachedUniversalIndexes[payload["key"]]
		if !ok {
			return fmt.Errorf("status bundle not found")
		}

		data, err := o.storage.ReadUniversalStatus(cursor)
		if err != nil {
			return err
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeGetStatusBundleResponse,
			Payload: json.RawMessage(fmt.Sprintf("%v", data)),
		}
		return o.swift.Send(ctx, response)
	})

	o.swift.RegisterHandler(swift.PacketTypeRawRequest, func(ctx context.Context, packet *swift.Packet) error {

		data, _ := structure.ParseOrderedMap(string(packet.Payload))
		signedRequest := model.NewSignedRequest(data)

		lastBlockHeight := storage.LastHeight()
		previousBlock, err := storage.GetChainStorageInstance().GetBlock(lastBlockHeight)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		// create a new machine for each request
		machine := vm.NewMachine(previousBlock)
		machine.Init(previousBlock)

		res, err := machine.Response(signedRequest)

		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		// Convert response to proper JSON format
		var jsonBytes []byte
		switch v := res.(type) {
		case string:
			jsonBytes = []byte(fmt.Sprintf(`"%s"`, v))
		case map[string]interface{}, []interface{}:
			// Handle nested JSON structures
			jsonBytes, err = json.Marshal(v)
			if err != nil {
				return o.swift.SendErrorResponse(ctx, "JSON marshaling error")
			}
		default:
			// For other types, use standard conversion
			jsonBytes = []byte(fmt.Sprintf(`%v`, v))
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeRawResponse,
			Payload: jsonBytes,
		}

		return o.swift.Send(ctx, response)
	})

	o.swift.RegisterHandler(swift.PacketTypeReplicateBlockRequest, func(ctx context.Context, packet *swift.Packet) error {
		o.mu.Lock()
		defer o.mu.Unlock()

		fmt.Println("replicate block request: ", string(packet.Payload))
		block, err := storage.ParseBlock(packet.Payload)

		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		err = o.machine.Commit(block)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		responseData, err := json.Marshal(map[string]string{"status": "success"})
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeReplicateBlockResponse,
			Payload: responseData,
		}

		return o.swift.Send(ctx, response)
	})

	o.swift.RegisterHandler(swift.PacketTypeRegisterReplicaRequest, func(ctx context.Context, packet *swift.Packet) error {
		responseData, err := json.Marshal(map[string]string{"status": "connected"})
		if err != nil {
			return err
		}

		var payload map[string]string
		err = json.Unmarshal(packet.Payload, &payload)
		if err != nil {
			return err
		}

		fmt.Println("payload: ", payload)
		o.RegisterReplica(payload["targetAddr"])

		response := &swift.Packet{
			Type:    swift.PacketTypeRegisterReplicaResponse,
			Payload: responseData,
		}
		return o.swift.Send(ctx, response)
	})

	// handshake request handler
	o.swift.RegisterHandler(swift.PacketTypeHandshakeCMDRequest, func(ctx context.Context, packet *swift.Packet) error {
		var payload struct {
			Peer string `json:"peer"`
		}

		swift.SwiftInfoLog("handshake request: %s", string(packet.Payload))
		if err := json.Unmarshal(packet.Payload, &payload); err != nil {
			return fmt.Errorf("handshake request payload parsing failed: %v", err)
		}

		if err := o.swift.Connect(payload.Peer); err != nil {
			return fmt.Errorf("failed to connect to peer: %v", err)
		}

		responseData, err := json.Marshal(map[string]string{"status": "connected"})
		if err != nil {
			return fmt.Errorf("failed to create response: %v", err)
		}

		swift.SwiftInfoLog("handshake response: %s", string(responseData))

		response := &swift.Packet{
			Type:    swift.PacketTypeHandshakeCMDResponse,
			Payload: responseData,
		}

		return o.swift.Send(ctx, response)
	})
}

func (o *Oracle) Shutdown() error {
	return util.TerminateProcess(storage.DataRootDir(), "oracle")
}

func (o *Oracle) RegisterReplica(targetAddr string) error {
	err := o.swift.Connect(targetAddr)
	if err != nil {
		return err
	}

	lastheight := storage.LastHeight()
	o.Replicas[targetAddr] = int64(lastheight)

	logger.Info("registered replica",
		zap.String("target_address", targetAddr),
	)

	fmt.Println("registered replica: ", o.Replicas)

	return nil
}

func (o *Oracle) BroadcastBlock(block *model.Block) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Skip if no replicas
	if len(o.Replicas) == 0 {
		return nil
	}

	// Prepare block data for replication
	blockData := block.Ser("full")

	// Create replication packet
	packet := &swift.Packet{
		Type:    swift.PacketTypeReplicateBlockRequest,
		Payload: json.RawMessage(blockData),
	}

	peers := make([]string, 0, len(o.Replicas))
	for peer := range o.Replicas {
		peers = append(peers, peer)
	}

	// Broadcast to all replicas
	logger.Info("broadcasting block to replicas",
		zap.Int("replica_count", len(o.Replicas)),
		zap.Int64("block_height", int64(block.Height)),
	)

	return o.swift.BroadcastPeers(peers, packet)
}
